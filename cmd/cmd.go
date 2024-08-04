package cmd

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/gabe565/telegram-to-mattermost/internal/etl"
	"github.com/gabe565/telegram-to-mattermost/internal/mattermost"
	"github.com/gabe565/telegram-to-mattermost/internal/telegram"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "telegram-to-mattermost dir",
		RunE: run,
		Args: cobra.ExactArgs(1),

		SilenceErrors: true,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("command missing context")
	}
	config.InitLog(false)
	conf.Input = args[0]

	export, err := telegram.FromFile(filepath.Join(conf.Input, "result.json"))
	if err != nil {
		return err
	}

	if err := etl.LoadUserMap(conf, export); err != nil {
		cmd.SilenceUsage = true
		return err
	}

	size, err := mattermost.TransformTelegramExport(conf, export)
	if err != nil {
		return err
	}

	slog.Info("Success!", "path", conf.Output, "size", humanize.IBytes(size))
	return nil
}
