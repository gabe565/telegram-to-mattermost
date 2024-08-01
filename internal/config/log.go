package config

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type clearLineWriter struct {
	*os.File
}

func (c clearLineWriter) Write(p []byte) (int, error) {
	_, err := c.File.Write(append([]byte("\r\x1B[K"), p...))
	return len(p), err
}

func InitLog(clearLine bool) {
	out := io.Writer(os.Stderr)
	if clearLine {
		out = clearLineWriter{os.Stderr}
	}
	slog.SetDefault(slog.New(log.NewWithOptions(out, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})))
}
