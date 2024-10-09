package util

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/require"
)

func TestSizeWriter(t *testing.T) {
	for _, i := range []int{humanize.Byte, humanize.KiByte, humanize.MiByte, humanize.GiByte} {
		w := SizeWriter{}
		_, err := w.Write(make([]byte, i))
		require.NoError(t, err)
		assert.Equal(t, uint64(i), w.Size()) //nolint:gosec
	}
}
