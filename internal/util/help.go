package util

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

func PrintPostRun(cmd *cobra.Command) {
	italic := lipgloss.NewStyle().Italic(true).Render
	_, _ = io.WriteString(
		cmd.OutOrStdout(),
		"\n"+
			lipgloss.JoinHorizontal(lipgloss.Top,
				"To import, run: ",
				italic("mmctl import upload data.zip")+" (or place in your data dir/bucket at "+italic("import/data.zip")+")\n"+
					italic("mmctl import list available")+"\n"+
					italic("mmctl import process ID_data.zip"),
			)+"\n"+
			"See Mattermost's "+
			termenv.Hyperlink("https://docs.mattermost.com/onboard/bulk-loading-data.html#bulk-load-data", "bulk loading data")+
			" guide for more info."+"\n",
	)
}
