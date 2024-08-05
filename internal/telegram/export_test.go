package telegram

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestExport_Users(t *testing.T) {
	numUsers := 7
	users := make([]User, 0, numUsers)
	for i := range numUsers {
		users = append(users, User{
			FromID: fmt.Sprintf("user%d", i),
			From:   fmt.Sprintf("User %d", i),
		})
	}

	numMessages := numUsers * 10
	messages := make([]*Message, 0, numMessages)
	for i := range numMessages {
		messages = append(messages, &Message{
			User: users[i%len(users)],
		})
	}

	export := Export{Messages: messages}
	got := export.Users()
	assert.Len(t, got, numUsers)
	for i, user := range got {
		assert.Equal(t, fmt.Sprintf("user%d", i), user.FromID)
		assert.Equal(t, fmt.Sprintf("User %d", i), user.From)
	}
	assert.Len(t, got, cap(got))
}

const unix = json.Number("1722806234")

//nolint:gochecknoglobals
var expect = time.Date(2024, time.August, 4, 21, 17, 14, 0, time.UTC)

func TestMessage_Date(t *testing.T) {
	m := Message{DateUnix: unix}
	assert.NotPanics(t, func() {
		assert.Equal(t, expect, m.Date().UTC())
	})
}

func TestMessage_Edited(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m Message
		assert.NotPanics(t, func() {
			assert.Nil(t, m.Edited())
		})
	})

	t.Run("not nil", func(t *testing.T) {
		m := Message{EditedUnix: ptr.To(unix)}
		assert.NotPanics(t, func() {
			got := m.Edited()
			require.NotNil(t, got)
			assert.Equal(t, expect, got.UTC())
		})
	})
}
