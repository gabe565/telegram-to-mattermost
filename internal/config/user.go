package config

type UserList map[string]*User

type User struct {
	TelegramUsername string `comment:"Telegram"`

	Username           string `comment:"Mattermost"`
	Email              string
	UseMarkdownPreview bool `default:"true"`
	UseFormatting      bool `default:"true"`
	ShowUnreadSection  bool `default:"false"`
}
