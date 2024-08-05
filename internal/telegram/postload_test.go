package telegram

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestExport_PostLoad(t *testing.T) {
	t.Run("replies", func(t *testing.T) {
		export := Export{Messages: []*Message{
			{ID: 1, DateUnix: "1"},
			{ID: 2, DateUnix: "2", ReplyToMessageID: ptr.To(int64(1))},
			{ID: 3, DateUnix: "3", ReplyToMessageID: ptr.To(int64(2))},
		}}
		assert.NotPanics(t, func() {
			export.PostLoad()
			require.Len(t, export.Messages, 1)
			require.NotNil(t, export.Messages[0].Reply)
			assert.EqualValues(t, 1, export.Messages[0].ID)
			assert.EqualValues(t, 2, export.Messages[0].Reply.ID)
			assert.EqualValues(t, 3, export.Messages[0].Reply.Reply.ID)
			assert.Len(t, export.Messages, cap(export.Messages))
		})
	})

	t.Run("pinned", func(t *testing.T) {
		export := Export{Messages: []*Message{
			{ID: 1, DateUnix: "1"},
			{ID: 2, DateUnix: "2", Event: &Event{
				Action:    ptr.To(ActionPinMessage),
				MessageID: ptr.To(int64(1)),
			}},
		}}
		assert.NotPanics(t, func() {
			export.PostLoad()
			require.Len(t, export.Messages, 1)
			assert.True(t, *export.Messages[0].IsPinned)
			assert.Len(t, export.Messages, cap(export.Messages))
		})
	})

	t.Run("events stripped", func(t *testing.T) {
		var export Export
		actions := []Action{
			ActionCreateGroup,
			ActionEditGroupPhoto,
			ActionInviteMembers,
			ActionJoinGroupByLink,
			ActionRemoveMembers,
			ActionPhoneCall,
			ActionScoreInGame,
			ActionEditChatTheme,
		}
		for i, action := range actions {
			export.Messages = append(export.Messages, &Message{
				ID:       int64(i),
				DateUnix: json.Number(strconv.Itoa(i)),
				Event:    &Event{Action: &action},
			})
		}
		assert.NotPanics(t, func() {
			export.PostLoad()
			require.Empty(t, export.Messages)
			assert.Len(t, export.Messages, cap(export.Messages))
		})
	})
}
