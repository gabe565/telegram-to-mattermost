package telegram

import (
	"encoding/json"
	"log/slog"
	"os"
)

func FromFile(path string, allowUnknown bool) (*Export, error) {
	slog.Info("Loading Telegram export", "path", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	decoder := json.NewDecoder(f)
	if !allowUnknown {
		decoder.DisallowUnknownFields()
	}
	export := &Export{}
	if err := decoder.Decode(export); err != nil {
		return nil, err
	}
	slog.Info("Loaded Telegram export", "name", export.Name, "type", export.Type, "messages", len(export.Messages))

	export.PostLoad()
	return export, nil
}
