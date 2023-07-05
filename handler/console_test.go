package handler_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestConsoleWithMaxLevel(t *testing.T) {
	l := slog.NewWithHandlers(handler.ConsoleWithMaxLevel(slog.InfoLevel))
	l.DoNothingOnPanicFatal()

	for _, level := range slog.AllLevels {
		l.Log(level, "a test message")
	}
	assert.NoErr(t, l.LastErr())
}
