package slog_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/formatter"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/processor"
)

func TestInfof(t *testing.T) {
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(formatter.NewJSONFormatter(slog.StringMap{
		"level": "levelName",
		"message": "msg",
		"data": "params",
	}))
	slog.AddHandler(h2)

	slog.AddProcessor(processor.AddHostname())

	slog.Infof("info %s", "message")
}
