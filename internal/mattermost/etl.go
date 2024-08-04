package mattermost

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/progressbar"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/gabe565/telegram-to-mattermost/internal/util"
	"github.com/mattermost/mattermost/server/v8/channels/app/imports"
	"github.com/mattermost/mattermost/server/v8/cmd/mmctl/commands/importer"
	"golang.org/x/image/webp"
)

func TransformTelegramExport(conf *config.Config, export *telegram.Export) (uint64, error) { //nolint:gocyclo
	slog.Info("Converting to Mattermost import", "path", conf.Output)

	f, err := os.Create(conf.Output)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = f.Close()
	}()

	var zw *zip.Writer
	var w io.Writer
	sizeWriter := &util.SizeWriter{}
	if strings.EqualFold(filepath.Ext(conf.Output), ".zip") {
		zw = zip.NewWriter(io.MultiWriter(f, sizeWriter))
		defer func() {
			_ = zw.Close()
		}()

		if w, err = zw.Create("data.jsonl"); err != nil {
			return 0, err
		}
	} else {
		w = io.MultiWriter(f, sizeWriter)
		if !conf.NoAttachments {
			slog.Warn(`Attachment paths will not be altered unless the output extension is ".zip"`)
		}
	}

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(Version()); err != nil {
		return 0, err
	}

	var channelType string
	if len(export.Users()) > 8 {
		channelType = importer.LineTypeChannel
	} else {
		channelType = importer.LineTypeDirectChannel
	}

	var team *imports.TeamImportData
	var channel *imports.ChannelImportData
	if channelType == importer.LineTypeChannel {
		teamLine := Team(export)
		team = teamLine.Team
		if err := encoder.Encode(teamLine); err != nil {
			return 0, err
		}

		channelLine := Channel(team)
		channel = channelLine.Channel
		if err := encoder.Encode(channelLine); err != nil {
			return 0, err
		}
	}

	if conf.CreateUsers {
		users := conf.Users
		slog.Info("Generating users", "count", len(users))
		for _, user := range users {
			if err := encoder.Encode(User(user, team)); err != nil {
				return 0, err
			}
		}
	}

	if channelType == importer.LineTypeDirectChannel {
		if err := encoder.Encode(DirectChannel(conf)); err != nil {
			return 0, err
		}
	}

	slog.Info("Generating posts", "count", len(export.Messages))
	bar := progressbar.New(len(export.Messages))
	var attachments []string
	for _, msg := range export.Messages {
		_ = bar.Add(1)
		if msg.From == "" {
			continue
		}

		var lines []imports.LineImportData
		switch channelType {
		case importer.LineTypeChannel:
			if lines, err = Post(conf, *team.Name, *channel.Name, msg); err != nil {
				return 0, err
			}
		case importer.LineTypeDirectChannel:
			if lines, err = DirectPost(conf, msg); err != nil {
				return 0, err
			}
		}

		for _, line := range lines {
			switch {
			case line.DirectPost != nil:
				if line.DirectPost.Attachments != nil {
					for _, attachment := range *line.DirectPost.Attachments {
						attachments = append(attachments, *attachment.Path)
					}
				}
				if line.DirectPost.Replies != nil {
					for _, msg := range *line.DirectPost.Replies {
						if msg.Attachments != nil {
							for _, attachment := range *msg.Attachments {
								attachments = append(attachments, *attachment.Path)
							}
						}
					}
				}
			case line.Post != nil:
				if line.Post.Attachments != nil {
					for _, attachment := range *line.Post.Attachments {
						attachments = append(attachments, *attachment.Path)
					}
				}
				if line.Post.Replies != nil {
					for _, msg := range *line.Post.Replies {
						if msg.Attachments != nil {
							for _, attachment := range *msg.Attachments {
								attachments = append(attachments, *attachment.Path)
							}
						}
					}
				}
			}

			if err := encoder.Encode(line); err != nil {
				return 0, err
			}
		}
	}
	_ = bar.Finish()

	if zw != nil {
		if !conf.NoAttachments {
			slog.Info("Zipping attachments (disable with --"+config.NoAttachmentsFlag+")", "count", len(attachments))
			bar = progressbar.New(len(attachments))
			for _, path := range attachments {
				_ = bar.Add(1)

				attachW, err := zw.Create(filepath.Join("data", path))
				if err != nil {
					return 0, err
				}

				if err := addAttachment(filepath.Join(conf.Input, path), attachW, !conf.NoFixWebP); err != nil {
					return 0, err
				}
			}
			_ = bar.Close()
		}

		if err := zw.Close(); err != nil {
			return 0, err
		}
	}

	return sizeWriter.Size(), f.Close()
}

var fixWebPWarned bool //nolint:gochecknoglobals

func addAttachment(path string, w io.Writer, fixWebP bool) error {
	src, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		path = strings.TrimSuffix(path, filepath.Ext(path)) + ".webp"
		if src, err = os.Open(path); err != nil {
			return err
		}
	}
	defer func() {
		_ = src.Close()
	}()

	if fixWebP && strings.EqualFold(filepath.Ext(path), ".webp") {
		if _, err := webp.Decode(src); err == nil {
			if _, err := src.Seek(0, io.SeekStart); err != nil {
				return err
			}
		} else {
			if !fixWebPWarned {
				fixWebPWarned = true
				slog.Warn("Detected unsupported WebP files. They will be converted using ImageMagick. (disable with --" + config.NoFixWebPFlag + ")")
			}
			fixed, err := util.BrokenWebPToPNG(path, src)
			defer func() {
				if fixed != nil {
					_ = fixed.Close()
				}
			}()
			if err != nil && !errors.Is(err, util.ErrNoImagemagick) {
				return err
			}

			_, err = io.Copy(w, fixed)
			return err
		}
	}

	_, err = io.Copy(w, src)
	return err
}
