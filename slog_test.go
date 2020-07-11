package slog_test

import (
	"os"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/processor"
)

func TestInfof(t *testing.T) {
	slog.AddHandler(&handler.StreamHandler{
		Out: os.Stdout,
	})
	slog.AddProcessor(processor.AddHostname())

	slog.Infof("info %s", "message")
}
