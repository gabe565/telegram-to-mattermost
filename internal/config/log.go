package config

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/schollz/progressbar/v3"
)

type clearLineWriter struct {
	bar *progressbar.ProgressBar
	*os.File
}

func (c clearLineWriter) Write(p []byte) (int, error) {
	_, err := c.File.Write(append([]byte("\r\x1B[K"), p...))
	if c.bar != nil {
		_, _ = c.File.WriteString(c.bar.String())
	}
	return len(p), err
}

func InitLog(bar *progressbar.ProgressBar) {
	out := io.Writer(os.Stderr)
	if bar != nil {
		out = clearLineWriter{bar, os.Stderr}
	}
	slog.SetDefault(slog.New(log.NewWithOptions(out, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})))
}
