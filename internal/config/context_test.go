package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.Background(), New())
	require.NotNil(t, ctx)
	conf, ok := ctx.Value(configKey).(*Config)
	assert.True(t, ok)
	assert.NotNil(t, conf)
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), configKey, New())
	conf, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, conf)
}
