package config

type Config struct {
	Input              string
	Output             string
	MaxTextLength      uint
	AllowUnknownFields bool

	CreateUsers bool
	UserFile    string
	Users       UserList

	NoAttachments bool
	NoFixWebP     bool
}

func New() *Config {
	return &Config{
		Output:        "data.zip",
		MaxTextLength: 4000,

		CreateUsers: true,
		UserFile:    "users.toml",
	}
}
