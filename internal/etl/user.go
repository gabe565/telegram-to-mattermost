package etl

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/mcuadros/go-defaults"
	"github.com/pelletier/go-toml/v2"
)

func LoadUserMap(conf *config.Config, export *telegram.Export) error {
	slog.Info("Loading user mapping", "path", conf.UserFile)

	allMappedUsers, err := loadMapping(conf.UserFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	exportUsers := export.Users()
	mappedUsers := make(config.UserList, len(exportUsers))
	for i, tgUser := range exportUsers {
		user, ok := allMappedUsers[tgUser.FromID]
		if !ok {
			user = &config.User{TelegramUsername: tgUser.From}
			defaults.SetDefaults(user)
			allMappedUsers[tgUser.FromID] = user
		}
		mappedUsers[tgUser.FromID] = user

		if user.Username == "" || (conf.CreateUsers && user.Email == "") {
			tbl := table.New().
				Row("ID:", tgUser.FromID).
				Row("Name:", user.TelegramUsername).
				Border(lipgloss.HiddenBorder()).
				StyleFunc(func(_, col int) lipgloss.Style {
					if col == 0 {
						return lipgloss.NewStyle().AlignHorizontal(lipgloss.Right)
					}
					return lipgloss.NewStyle()
				})

			if err := huh.NewForm(
				huh.NewGroup(
					huh.NewNote().
						Title(fmt.Sprintf("Map Telegram user (%d/%d)\n%s", i+1, len(exportUsers), tbl.String())),

					huh.NewInput().
						Title("Username").
						Validate(huh.ValidateNotEmpty()).
						Value(&user.Username),

					huh.NewInput().
						Title("Email").
						Validate(huh.ValidateNotEmpty()).
						Value(&user.Email),
				),
			).Run(); err != nil {
				return err
			}

			if err := saveMapping(conf.UserFile, allMappedUsers); err != nil {
				return err
			}
		}
	}

	channelMembers := make([]string, 0, len(mappedUsers))
	for _, u := range mappedUsers {
		channelMembers = append(channelMembers, u.Username)
	}
	conf.ChannelMembers = &channelMembers
	conf.Users = mappedUsers
	return nil
}

func loadMapping(path string) (config.UserList, error) {
	mapping := make(config.UserList)

	f, err := os.Open(path)
	if err != nil {
		return mapping, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := toml.NewDecoder(f).Decode(&mapping); err != nil {
		return mapping, err
	}

	return mapping, nil
}

func saveMapping(path string, mapping config.UserList) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	encoder := toml.NewEncoder(f)
	encoder.SetIndentTables(true)
	if err := encoder.Encode(mapping); err != nil {
		return err
	}

	return f.Close()
}
