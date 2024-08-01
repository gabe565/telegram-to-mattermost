package util

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

var ErrNoImagemagick = errors.New("imagemagick is not installed")

type RWSeekCloser interface {
	io.ReadSeekCloser
	io.Writer
	Name() string
}

func IsBrokenWebP(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := webp.Decode(f); err == nil {
		return false, nil
	}
	return true, nil
}

func BrokenWebPToPNG(path string, f RWSeekCloser) (RWSeekCloser, error) {
	if _, err := webp.Decode(f); err == nil {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return f, nil
	}

	if _, err := exec.LookPath("magick"); err != nil {
		slog.Warn("Found an unsupported webp. Please install ImageMagick to fix.", "path", path)
		return nil, ErrNoImagemagick
	}

	_ = f.Close()

	temp := filepath.Join(os.TempDir(), filepath.Base(path)+".png")

	var errBuf strings.Builder

	cmd := exec.Command("magick", path, temp)
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, errBuf.String())
	}

	f, err := os.Open(temp)
	if err != nil {
		return nil, err
	}

	return TempFile{RWSeekCloser: f}, nil
}

type TempFile struct {
	RWSeekCloser
}

func (t TempFile) Close() error {
	return errors.Join(
		t.RWSeekCloser.Close(),
		os.Remove(t.RWSeekCloser.Name()),
	)
}
