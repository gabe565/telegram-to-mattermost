package util

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/require"
)

func TestSizeWriter(t *testing.T) {
	for _, i := range []int{1, 1024, 1024 * 1024, 1024 * 1024 * 1024} {
		w := SizeWriter{}
		_, err := w.Write(make([]byte, i))
		require.NoError(t, err)
		assert.Equal(t, int64(i), w.Size())
	}
}
