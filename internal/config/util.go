package config

import (
	"errors"
	"fmt"
)

var ErrUserNotInMapping = errors.New("user is not in mapping")

func (c *Config) MapUser(username string) (string, error) {
	mapped, ok := c.Usernames[username]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUserNotInMapping, username)
	}
	return mapped, nil
}

var ErrEmailNotInMapping = errors.New("email is not in mapping")

func (c *Config) MapEmail(email string) (string, error) {
	mapped, ok := c.Emails[email]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrEmailNotInMapping, email)
	}
	return mapped, nil
}
