package etl

import (
	"errors"
	"fmt"
	"os"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/mcuadros/go-defaults"
	"github.com/pelletier/go-toml/v2"
)

var ErrNotAllMapped = errors.New("not all users are mapped")

func LoadUserMap(conf *config.Config, export *telegram.Export) error {
	exportUsers := export.Users()
	allMappedUsers := make(config.UserList, len(exportUsers))

	if _, err := os.Stat(conf.UserFile); err == nil {
		f, err := os.Open(conf.UserFile)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		if err := toml.NewDecoder(f).Decode(&allMappedUsers); err != nil {
			return err
		}

		_ = f.Close()
	}

	mappedUsers := make(config.UserList, len(exportUsers))
	var missing bool
	for _, tgUser := range exportUsers {
		user, ok := allMappedUsers[tgUser.FromID]
		if !ok {
			missing = true
			user = &config.User{
				TelegramUsername: tgUser.From,
			}
			defaults.SetDefaults(user)
			allMappedUsers[tgUser.FromID] = user
		} else if !missing {
			missing = user.Username == "" || (conf.CreateUsers && user.Email == "")
			mappedUsers[tgUser.FromID] = user
		}
	}

	if missing {
		f, err := os.Create(conf.UserFile)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		encoder := toml.NewEncoder(f)
		encoder.SetIndentTables(true)
		if err := encoder.Encode(allMappedUsers); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		return fmt.Errorf("%w: Please edit %q with mapping details, then rerun this tool", ErrNotAllMapped, conf.UserFile)
	}

	conf.Users = mappedUsers
	return nil
}
