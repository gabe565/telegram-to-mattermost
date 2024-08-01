package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/telegram-to-mattermost/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("204"))
		fmt.Println(lipgloss.JoinHorizontal(
			lipgloss.Top, errStyle.Render("Error: "), err.Error(),
		))
		os.Exit(1)
	}
}
