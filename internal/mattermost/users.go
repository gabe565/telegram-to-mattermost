package mattermost

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
)

var ErrUserNotMapped = errors.New("usernames not mapped")

func MapAllUsers(conf *config.Config, export *telegram.Export) ([]string, error) {
	users := export.Users()
	var missing []string
	for k, v := range users {
		var ok bool
		users[k], ok = conf.Usernames[v]
		if !ok {
			missing = append(missing, v)
		}
	}
	if len(missing) != 0 {
		return nil, fmt.Errorf("%w: %s", ErrUserNotMapped, strings.Join(missing, ", "))
	}
	return users, nil
}

var ErrEmailNotMapped = errors.New("emails not mapped")

func CheckEmails(conf *config.Config, export *telegram.Export) error {
	var missing []string
	for _, v := range export.Users() {
		if _, ok := conf.Emails[v]; !ok {
			missing = append(missing, v)
		}
	}
	if len(missing) != 0 {
		return fmt.Errorf("%w: %s", ErrEmailNotMapped, strings.Join(missing, ", "))
	}
	return nil
}
