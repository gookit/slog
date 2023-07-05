package slog_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
)

func TestNewLvFormatter(t *testing.T) {
	lf := slog.NewLvFormatter(slog.InfoLevel)

	assert.True(t, lf.IsHandling(slog.ErrorLevel))
	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.False(t, lf.IsHandling(slog.DebugLevel))

	lf.SetMaxLevel(slog.DebugLevel)
	assert.True(t, lf.IsHandling(slog.DebugLevel))
}

func TestNewLvsFormatter(t *testing.T) {
	lf := slog.NewLvsFormatter([]slog.Level{slog.InfoLevel, slog.ErrorLevel})
	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.False(t, lf.IsHandling(slog.DebugLevel))

	lf.SetLimitLevels([]slog.Level{slog.InfoLevel, slog.ErrorLevel, slog.DebugLevel})
	assert.True(t, lf.IsHandling(slog.DebugLevel))
}

func TestLevelFormatting(t *testing.T) {
	lf := slog.NewMaxLevelFormatting(slog.InfoLevel)

	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.False(t, lf.IsHandling(slog.TraceLevel))

	// use levels
	lf = slog.NewLevelsFormatting([]slog.Level{slog.InfoLevel, slog.ErrorLevel})

	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.True(t, lf.IsHandling(slog.ErrorLevel))
	assert.False(t, lf.IsHandling(slog.TraceLevel))
}
