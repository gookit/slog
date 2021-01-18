package slog_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

func TestLevel_Name(t *testing.T) {
	assert.Equal(t, "INFO", slog.InfoLevel.Name())
	assert.Equal(t, "INFO", slog.InfoLevel.String())
	assert.Equal(t, "info", slog.InfoLevel.LowerName())
}

func TestLevel_ShouldHandling(t *testing.T) {
	assert.True(t, slog.InfoLevel.ShouldHandling(slog.ErrorLevel))
	assert.False(t, slog.InfoLevel.ShouldHandling(slog.TraceLevel))

	assert.True(t, slog.DebugLevel.ShouldHandling(slog.InfoLevel))
	assert.False(t, slog.DebugLevel.ShouldHandling(slog.TraceLevel))
}
