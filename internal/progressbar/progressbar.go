package progressbar

import (
	"io"
	"os"
	"time"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/schollz/progressbar/v3"
)

type Align uint8

const (
	Left Align = iota
	Right
)

func New(total int, alignDescription Align) *progressbar.ProgressBar {
	opts := []progressbar.Option{
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(65 * time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			_, _ = io.WriteString(os.Stderr, "\n")
			config.InitLog(nil)
		}),
	}
	if alignDescription == Right {
		opts = append(opts, progressbar.OptionShowDescriptionAtLineEnd())
	}
	bar := progressbar.NewOptions(total, opts...)
	config.InitLog(bar)
	return bar
}
