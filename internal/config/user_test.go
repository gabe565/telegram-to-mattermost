package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserList_Usernames(t *testing.T) {
	tests := []struct {
		name string
		u    UserList
		want []string
	}{
		{"empty", UserList{}, []string{}},
		{"DM", UserList{"user1": {Username: "gabe565"}, "user2": {Username: "2"}}, []string{"2", "gabe565"}},
		{"saved messages", UserList{"user1": {Username: "gabe565"}}, []string{"gabe565", "gabe565"}},
		{"group", UserList{"user1": {Username: "gabe565"}, "user2": {Username: "2"}, "user3": {Username: "3"}}, []string{"2", "3", "gabe565"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.Usernames())
		})
	}
}
