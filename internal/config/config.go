package config

type Config struct {
	Input         string
	Output        string
	MaxTextLength uint

	CreateUsers bool
	UserFile    string
	Users       UserList

	NoAttachments bool
	NoFixWebP     bool

	ChannelMembers *[]string
}

func New() *Config {
	return &Config{
		Output:        "data.zip",
		MaxTextLength: 4000,
		UserFile:      "users.toml",
	}
}
