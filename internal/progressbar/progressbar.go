package progressbar

import (
	"io"
	"os"
	"time"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/schollz/progressbar/v3"
)

func New(max int) *progressbar.ProgressBar {
	config.InitLog(true)
	return progressbar.NewOptions(max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			_, _ = io.WriteString(os.Stderr, "\n")
			config.InitLog(false)
		}),
	)
}
