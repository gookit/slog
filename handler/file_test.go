package handler_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestNewFileHandler(t *testing.T) {
	h := handler.NewFileHandler("./testdata/app.log", false)

	l := slog.NewWithHandlers(h)
	l.Info("message")
}
