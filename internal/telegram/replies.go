package telegram

import (
	"log/slog"
	"slices"

	"github.com/gabe565/telegram-to-mattermost/internal/progressbar"
)

func (e *Export) GroupReplies() {
	slog.Info("Grouping replies")
	bar := progressbar.New(len(e.Messages))
	var deleteIDs []int64
	var replyCount int
	for _, msg := range e.Messages {
		_ = bar.Add(1)
		if msg.ReplyToMessageID != nil {
			for _, replied := range e.Messages {
				if *msg.ReplyToMessageID == replied.ID {
					replied.Reply = msg
					deleteIDs = append(deleteIDs, msg.ID)
					replyCount++
					break
				}
			}
		}
	}
	_ = bar.Close()
	slog.Info("Done grouping replies", "count", replyCount)

	e.Messages = slices.DeleteFunc(e.Messages, func(msg *Message) bool {
		return slices.Contains(deleteIDs, msg.ID)
	})
	e.Messages = slices.Clip(e.Messages)
}
