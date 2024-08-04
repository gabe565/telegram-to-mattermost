package telegram

import (
	"log/slog"
	"slices"

	"github.com/gabe565/telegram-to-mattermost/internal/progressbar"
	"k8s.io/utils/ptr"
)

func (e *Export) PostLoad() {
	slog.Info("Running post-load actions")
	bar := progressbar.New(len(e.Messages))

	var deleteIDs []int64
	var replyCount, pinCount int

	for _, msg := range e.Messages {
		switch {
		case msg.Event != nil && msg.Event.Action != nil:
			switch *msg.Event.Action {
			case ActionPinMessage:
				for _, pinned := range e.Messages {
					if *msg.Event.MessageID == pinned.ID {
						pinned.IsPinned = ptr.To(true)
						deleteIDs = append(deleteIDs, msg.ID)
						pinCount++
						break
					}
				}
			default:
				deleteIDs = append(deleteIDs, msg.ID)
			}
		case msg.ReplyToMessageID != nil:
			for _, replied := range e.Messages {
				if *msg.ReplyToMessageID == replied.ID {
					replied.Reply = msg
					deleteIDs = append(deleteIDs, msg.ID)
					replyCount++
					break
				}
			}
		}
		_ = bar.Add(1)
	}
	_ = bar.Close()
	slog.Info("Finished post-load actions", slog.Group("counts", "replies", replyCount, "pins", pinCount))

	e.Messages = slices.DeleteFunc(e.Messages, func(msg *Message) bool {
		return slices.Contains(deleteIDs, msg.ID)
	})
	e.Messages = slices.Clip(e.Messages)
}
