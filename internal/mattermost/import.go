package mattermost

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/gabe565/telegram-to-mattermost/internal/util"
	"github.com/mattermost/mattermost/server/v8/channels/app/imports"
	"github.com/mattermost/mattermost/server/v8/cmd/mmctl/commands/importer"
	"k8s.io/utils/ptr"
)

func Version() *imports.LineImportData {
	return &imports.LineImportData{
		Type:    importer.LineTypeVersion,
		Version: ptr.To(1),
		Info: &imports.VersionInfoImportData{
			Generator: "gabe565/telegram-to-mattermost",
			Version:   "", // TODO: Pass version
			Created:   time.Now().Format(time.RFC3339Nano),
		},
	}
}

func User(user *config.User) *imports.LineImportData {
	userImport := &imports.UserImportData{
		Username:           &user.Username,
		Email:              &user.Username,
		UseMarkdownPreview: ptr.To(strconv.FormatBool(user.UseMarkdownPreview)),
		UseFormatting:      ptr.To(strconv.FormatBool(user.UseFormatting)),
		ShowUnreadSection:  ptr.To(strconv.FormatBool(user.ShowUnreadSection)),
		EmailInterval:      ptr.To("immediately"),
	}

	return &imports.LineImportData{Type: importer.LineTypeUser, User: userImport}
}

func DirectChannel(conf *config.Config) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeDirectChannel,
		DirectChannel: &imports.DirectChannelImportData{
			Members: conf.ChannelMembers,
		},
	}
}

func DirectPost(conf *config.Config, msg *telegram.Message) (*imports.LineImportData, error) {
	createAt, err := transformTimestamp(&msg.Date)
	if err != nil {
		return nil, err
	}

	editAt, err := transformTimestamp(msg.Edited)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]
	post := &imports.DirectPostImportData{
		ChannelMembers: conf.ChannelMembers,
		User:           &user.Username,
		Message:        ptr.To(msg.FormatText(conf.MaxTextLength)),
		CreateAt:       createAt,
		EditAt:         editAt,
	}

	if attachment, err := transformAttachment(conf, msg); err != nil {
		return nil, err
	} else if attachment != nil {
		post.Attachments = &[]imports.AttachmentImportData{*attachment}
	}

	for msg := msg.Reply; msg != nil; msg = msg.Reply {
		replyImport, err := Reply(conf, msg)
		if err != nil {
			return nil, err
		}

		if post.Replies == nil {
			post.Replies = ptr.To(make([]imports.ReplyImportData, 0, 1))
		}
		*post.Replies = append(*post.Replies, *replyImport)
	}

	if msg.IsPinned != nil && *msg.IsPinned == true {
		post.IsPinned = ptr.To(true)
	}

	return &imports.LineImportData{
		Type:       importer.LineTypeDirectPost,
		DirectPost: post,
	}, nil
}

func Reply(conf *config.Config, msg *telegram.Message) (*imports.ReplyImportData, error) {
	createAt, err := transformTimestamp(&msg.Date)
	if err != nil {
		return nil, err
	}

	editAt, err := transformTimestamp(msg.Edited)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]
	post := &imports.ReplyImportData{
		User:     &user.Username,
		Message:  ptr.To(msg.FormatText(conf.MaxTextLength)),
		CreateAt: createAt,
		EditAt:   editAt,
	}

	if attachment, err := transformAttachment(conf, msg); err != nil {
		return nil, err
	} else if attachment != nil {
		post.Attachments = &[]imports.AttachmentImportData{*attachment}
	}

	return post, nil
}

func transformTimestamp(v *json.Number) (*int64, error) {
	if v == nil {
		return nil, nil
	}

	parsed, err := v.Int64()
	if err != nil {
		return nil, err
	}
	parsed *= 1000
	return &parsed, nil
}

func transformAttachment(conf *config.Config, msg *telegram.Message) (*imports.AttachmentImportData, error) {
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

	return &imports.AttachmentImportData{Path: &path}, nil //dst.Close()
}
