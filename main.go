package main

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/telegram-to-mattermost/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("204"))
		_, _ = io.WriteString(os.Stderr, lipgloss.JoinHorizontal(
			lipgloss.Top, errStyle.Render("Error: "), err.Error(),
		)+"\n")
		os.Exit(1)
	}
}
