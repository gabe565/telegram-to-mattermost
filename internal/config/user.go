package config

import (
	"maps"
	"slices"
)

type UserList map[string]*User

func (u UserList) Usernames() []string {
	result := make([]string, 0, len(u))
	for _, user := range u {
		result = append(result, user.Username)
	}
	if len(result) == 1 {
		result = append(result, result[0])
	} else {
		slices.Sort(result)
		result = slices.Compact(result)
	}
	return result
}

func (u UserList) Unique() UserList {
	result := maps.Clone(u)
	seen := make([]string, 0, len(result))
	for k, v := range result {
		if slices.Contains(seen, v.Username) {
			delete(result, k)
		} else {
			seen = append(seen, v.Username)
		}
	}
	return result
}

type User struct {
	TelegramUsername string `comment:"Telegram"`

	Username           string `comment:"Mattermost"`
	Email              string
	UseMarkdownPreview bool `default:"true"`
	UseFormatting      bool `default:"true"`
	ShowUnreadSection  bool `default:"false"`
}
