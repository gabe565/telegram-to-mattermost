package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	conf := New()
	require.NotNil(t, conf)
}
