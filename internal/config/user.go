package config

type UserList map[string]*User

func (u UserList) Usernames() []string {
	result := make([]string, 0, len(u))
	for _, user := range u {
		result = append(result, user.Username)
	}
	if len(result) == 1 {
		result = append(result, result[0])
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
