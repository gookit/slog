package slog_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
)

func TestSafeToLevelMode(t *testing.T) {
	assert.Eq(t, slog.LevelModeList, slog.SafeToLevelMode("list"))
	assert.Eq(t, slog.LevelModeList, slog.SafeToLevelMode("0"))
	assert.Eq(t, slog.LevelModeMax, slog.SafeToLevelMode("1"))
	assert.Eq(t, slog.LevelModeList, slog.SafeToLevelMode("unknown"))

	mode := slog.SafeToLevelMode("max")
	assert.Eq(t, slog.LevelModeMax, mode)

	// MarshalJSON
	bs, err := mode.MarshalJSON()
	assert.Nil(t, err)
	assert.Eq(t, `"max"`, string(bs))

	// UnmarshalJSON
	mode = slog.LevelMode(0)
	err = mode.UnmarshalJSON([]byte(`"max"`))
	assert.Nil(t, err)
	assert.Eq(t, slog.LevelModeMax, mode)

	assert.Err(t, mode.UnmarshalJSON([]byte("ab")))
}

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

	// test level mode
	assert.Eq(t, "list", slog.LevelModeList.String())
	assert.Eq(t, "max", slog.LevelModeMax.String())
	assert.Eq(t, "unknown", slog.LevelMode(9).String())
}
