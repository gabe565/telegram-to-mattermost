package mattermost

import (
	"os"
	"path/filepath"
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

func User(user *config.User, team *imports.TeamImportData) *imports.LineImportData {
	userImport := &imports.UserImportData{
		Username:           &user.Username,
		Email:              &user.Username,
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
		})
	}

	return &imports.LineImportData{Type: importer.LineTypeUser, User: userImport}
}

func Team(export *telegram.Export) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeTeam,
		Team: &imports.TeamImportData{
			Name:        ptr.To(xstrings.ToKebabCase(export.Name)),
			DisplayName: ptr.To(export.Name),
			Type:        ptr.To("I"),
		},
	}
}

func Channel(team *imports.TeamImportData) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeChannel,
		Channel: &imports.ChannelImportData{
			Team:        team.Name,
			Name:        team.Name,
			DisplayName: team.DisplayName,
			Type:        ptr.To(model.ChannelTypePrivate),
		},
	}
}

func DirectChannel(conf *config.Config) *imports.LineImportData {
	return &imports.LineImportData{
		Type: importer.LineTypeDirectChannel,
		DirectChannel: &imports.DirectChannelImportData{
			Members: conf.ChannelMembers,
		},
	}
}

func Post(conf *config.Config, team, channel string, msg *telegram.Message) (*imports.LineImportData, error) {
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

	post := &imports.PostImportData{
		Team:        &team,
		Channel:     &channel,
		User:        &user.Username,
		Message:     ptr.To(msg.FormatText(conf)),
		CreateAt:    createAt,
		EditAt:      editAt,
		Replies:     replies,
		Attachments: attachments,
	}

	if msg.IsPinned != nil && *msg.IsPinned == true {
		post.IsPinned = ptr.To(true)
	}

	return &imports.LineImportData{
		Type: importer.LineTypePost,
		Post: post,
	}, nil
}

func DirectPost(conf *config.Config, msg *telegram.Message) (*imports.LineImportData, error) {
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

	post := &imports.DirectPostImportData{
		ChannelMembers: conf.ChannelMembers,
		User:           &user.Username,
		Message:        ptr.To(msg.FormatText(conf)),
		CreateAt:       createAt,
		EditAt:         editAt,
		Replies:        replies,
		Attachments:    attachments,
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
	createAt, editAt := timestamps(msg)

	attachments, err := transformAttachment(conf, msg)
	if err != nil {
		return nil, err
	}

	user := conf.Users[msg.FromID]

	return &imports.ReplyImportData{
		User:        &user.Username,
		Message:     ptr.To(msg.FormatText(conf)),
		CreateAt:    createAt,
		EditAt:      editAt,
		Attachments: attachments,
	}, nil
}

func timestamps(msg *telegram.Message) (*int64, *int64) {
	createAt := msg.Date().UnixMilli()
	var editAt *int64
	if editedDate := msg.Edited(); editedDate != nil {
		editAt = ptr.To(editedDate.UnixMilli())
	}
	return &createAt, editAt
}

func transformReplies(conf *config.Config, msg *telegram.Message) (*[]imports.ReplyImportData, error) {
	var replies []imports.ReplyImportData
	for msg := msg.Reply; msg != nil; msg = msg.Reply {
		replyImport, err := Reply(conf, msg)
		if err != nil {
			return nil, err
		}

		replies = append(replies, *replyImport)
	}
	if len(replies) == 0 {
		return nil, nil
	}
	return &replies, nil
}

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
