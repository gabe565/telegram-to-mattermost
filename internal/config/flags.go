package config

import "github.com/spf13/cobra"

const (
	OutputFlag        = "output"
	MaxTextLengthFlag = "max-text-length"
	CreateUsersFlag   = "create-users"
	MapUsernamesFlag  = "map-usernames"
	EmailsFlag        = "map-emails"
	NoAttachmentsFlag = "no-attachments"
	NoFixWebPFlag     = "no-fix-webp"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.StringVarP(&c.Output, OutputFlag, "o", c.Output, "Output filename")
	fs.UintVar(&c.MaxTextLength, MaxTextLengthFlag, c.MaxTextLength, "Maximum post text length")
	fs.BoolVar(&c.CreateUsers, CreateUsersFlag, c.CreateUsers, "Adds users to Mattermost import")
	fs.StringToStringVar(&c.Usernames, MapUsernamesFlag, c.Usernames, "Map of Telegram usernames to Mattermost usernames")
	fs.StringToStringVar(&c.Emails, EmailsFlag, c.Emails, "Map of Telegram usernames to Mattermost emails")
	fs.BoolVar(&c.NoAttachments, NoAttachmentsFlag, c.NoAttachments, "Disables embedding of attachments")
	fs.BoolVar(&c.NoFixWebP, NoFixWebPFlag, c.NoFixWebP, "Disable fixing of WebP files (usually stickers) which will not load into Mattermost")
}
