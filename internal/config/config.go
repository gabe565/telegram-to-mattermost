package config

type Config struct {
	Input         string
	Output        string
	MaxTextLength uint

	CreateUsers bool
	Usernames   map[string]string
	Emails      map[string]string

	NoAttachments bool
	NoFixWebP     bool

	ChannelMembers *[]string
}

func New() *Config {
	return &Config{
		Output:        "data.zip",
		MaxTextLength: 4000,
	}
}
