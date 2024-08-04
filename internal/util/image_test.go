package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTempFile_Close(t *testing.T) {
	temp, err := os.CreateTemp("", "tg-to-mm")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = temp.Close()
		_ = os.Remove(temp.Name())
	})

	tempFile := TempFile{temp}
	require.NoError(t, tempFile.Close())
	_, err = os.Stat(tempFile.Name())
	assert.True(t, os.IsNotExist(err))
}
