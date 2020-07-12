package slog_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/processor"
)

func TestInfof(t *testing.T) {
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))
	slog.AddProcessor(processor.AddHostname())

	slog.Infof("info %s", "message")
}
