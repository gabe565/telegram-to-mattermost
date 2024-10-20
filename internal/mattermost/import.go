package mattermost

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/gabe565/telegram-to-mattermost/internal/util"
	"github.com/huandu/xstrings"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/imports"
	"github.com/mattermost/mattermost/server/v8/cmd/mmctl/commands/importer"
	"k8s.io/utils/ptr"
)

func Version() *imports.LineImportData {
	var commit string
	var modified bool
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}
	}
	if len(commit) > 8 {
		commit = commit[:8]
	}
	if modified {
		commit = "*" + commit
	}

	return &imports.LineImportData{
		Type:    importer.LineTypeVersion,
		Version: ptr.To(1),
		Info: &imports.VersionInfoImportData{
			Generator: "gabe565/telegram-to-mattermost",
			Version:   commit,
			Created:   time.Now().Format(time.RFC3339Nano),
		},
	}
}

func User(team *imports.TeamImportData, channel *imports.ChannelImportData, user *config.User) *imports.LineImportData {
	userImport := &imports.UserImportData{
		Username:           &user.Username,
		Email:              &user.Email,
		UseMarkdownPreview: ptr.To(strconv.FormatBool(user.UseMarkdownPreview)),
		UseFormatting:      ptr.To(strconv.FormatBool(user.UseFormatting)),
		ShowUnreadSection:  ptr.To(strconv.FormatBool(user.ShowUnreadSection)),
		EmailInterval:      ptr.To("immediately"),
	}

	if team != nil {
		if userImport.Teams == nil {
			userImport.Teams = ptr.To(make([]imports.UserTeamImportData, 0, 1))
		}
		*userImport.Teams = append(*userImport.Teams, imports.UserTeamImportData{
			Name:  team.Name,
			Roles: ptr.To("team_user"),
			Channels: &[]imports.UserChannelImportData{{
				Name:  channel.Name,
				Roles: ptr.To(model.ChannelUserRoleId),
			}},
		})
	}

	return &imports.LineImportData{Type: importer.LineTypeUser, User: userImport}
}

func Team(conf *config.Config) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeTeam,
		Team: &imports.TeamImportData{
			Name:        ptr.To(xstrings.ToKebabCase(conf.TeamName)),
			DisplayName: &conf.TeamName,
			Type:        ptr.To("I"),
		},
	}
}

func Channel(export *telegram.Export, team *imports.TeamImportData) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeChannel,
		Channel: &imports.ChannelImportData{
			Team:        team.Name,
			Name:        ptr.To(xstrings.ToKebabCase(export.Name)),
			DisplayName: &export.Name,
			Type:        ptr.To(model.ChannelTypePrivate),
		},
	}
}

func DirectChannel(conf *config.Config) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeDirectChannel,
		DirectChannel: &imports.DirectChannelImportData{
			Members: ptr.To(conf.Users.Usernames()),
		},
	}
}

func Post(conf *config.Config, team, channel string, msg *telegram.Message) ([]imports.LineImportData, error) {
	createAt, editAt := timestamps(msg)

	attachments, err := transformAttachment(conf, msg)
	if err != nil {
		return nil, err
	}

	replies, err := transformReplies(conf, msg)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]

	texts := msg.FormatText(conf)
	lines := make([]imports.LineImportData, 0, len(texts))
	for _, text := range texts {
		post := &imports.PostImportData{
			Team:        &team,
			Channel:     &channel,
			User:        &user.Username,
			Message:     &text,
			CreateAt:    createAt,
			EditAt:      editAt,
			Replies:     replies,
			Attachments: attachments,
		}

		if msg.IsPinned != nil && *msg.IsPinned {
			post.IsPinned = ptr.To(true)
		}

		lines = append(lines, imports.LineImportData{
			Type: importer.LineTypePost,
			Post: post,
		})
	}

	return lines, nil
}

func DirectPost(conf *config.Config, msg *telegram.Message) ([]imports.LineImportData, error) {
	createAt, editAt := timestamps(msg)

	attachments, err := transformAttachment(conf, msg)
	if err != nil {
		return nil, err
	}

	replies, err := transformReplies(conf, msg)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]

	texts := msg.FormatText(conf)
	lines := make([]imports.LineImportData, 0, len(texts))
	for _, text := range texts {
		post := &imports.DirectPostImportData{
			ChannelMembers: ptr.To(conf.Users.Usernames()),
			User:           &user.Username,
			Message:        &text,
			CreateAt:       createAt,
			EditAt:         editAt,
			Replies:        replies,
			Attachments:    attachments,
		}

		if msg.IsPinned != nil && *msg.IsPinned {
			post.IsPinned = ptr.To(true)
		}

		lines = append(lines, imports.LineImportData{
			Type:       importer.LineTypeDirectPost,
			DirectPost: post,
		})
	}

	return lines, nil
}

func Reply(conf *config.Config, msg *telegram.Message) ([]imports.ReplyImportData, error) {
	createAt, editAt := timestamps(msg)

	attachments, err := transformAttachment(conf, msg)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]

	texts := msg.FormatText(conf)
	lines := make([]imports.ReplyImportData, 0, len(texts))
	for _, text := range texts {
		lines = append(lines, imports.ReplyImportData{
			User:        &user.Username,
			Message:     &text,
			CreateAt:    createAt,
			EditAt:      editAt,
			Attachments: attachments,
		})
	}
	return lines, nil
}

func timestamps(msg *telegram.Message) (*int64, *int64) {
	createAt := msg.Date().UnixMilli()
	var editAt *int64
	if editedDate := msg.Edited(); editedDate != nil {
		editAt = ptr.To(editedDate.UnixMilli())
	}
	return &createAt, editAt
}

//nolint:nilnil
func transformReplies(conf *config.Config, msg *telegram.Message) (*[]imports.ReplyImportData, error) {
	var replies []imports.ReplyImportData
	for msg := msg.Reply; msg != nil; msg = msg.Reply {
		replyImports, err := Reply(conf, msg)
		if err != nil {
			return nil, err
		}

		replies = append(replies, replyImports...)
	}
	if len(replies) == 0 {
		return nil, nil
	}
	return &replies, nil
}

//nolint:nilnil
func transformAttachment(conf *config.Config, msg *telegram.Message) (*[]imports.AttachmentImportData, error) {
	if conf.NoAttachments {
		return nil, nil
	}

	var path string
	switch {
	case msg.File != nil && msg.File.Path != nil:
		path = *msg.File.Path
	case msg.File != nil && msg.File.Photo != nil:
		path = *msg.File.Photo
	case msg.Contact != nil && msg.Contact.VCard != nil:
		path = *msg.Contact.VCard
	default:
		return nil, nil
	}

	if _, err := os.Stat(filepath.Join(conf.Input, path)); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if !conf.NoFixWebP && strings.EqualFold(filepath.Ext(path), ".webp") {
		broken, err := util.IsBrokenWebP(filepath.Join(conf.Input, path))
		if err != nil {
			return nil, err
		}
		if broken {
			path = strings.TrimSuffix(path, filepath.Ext(path)) + ".png"
		}
	}

	return &[]imports.AttachmentImportData{{Path: &path}}, nil
}
