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
				Row("ID", tgUser.FromID).
				Row("Name", user.TelegramUsername).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
				StyleFunc(func(_, col int) lipgloss.Style {
					s := lipgloss.NewStyle().Padding(0, 1)
					if col == 0 {
						s = s.AlignHorizontal(lipgloss.Right)
					}
					return s
				})

			if err := huh.NewForm(
				huh.NewGroup(
					huh.NewNote().
						Title(fmt.Sprintf("Map Telegram User (%d/%d)", i+1, len(exportUsers))).
						Description(tbl.String()),

					huh.NewInput().
						Title("Username").
						Validate(huh.ValidateNotEmpty()).
						Value(&user.Username),

					huh.NewInput().
						Title("Email").
						Validate(huh.ValidateNotEmpty()).
						Value(&user.Email),
				),
			).WithTheme(huh.ThemeDracula()).Run(); err != nil {
				return err
			}

			if err := saveMapping(conf.UserFile, allMappedUsers); err != nil {
				return err
			}
		}
	}

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
