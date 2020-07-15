package slog_test

import (
	"fmt"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestNewLineFormatter(t *testing.T) {
	lf := slog.NewLineFormatter()

	fmt.Println(lf.FieldMap())
}

func TestJSONFormatter(t *testing.T) {
	l := slog.New()

	f := slog.NewJSONFormatter(nil)
	h := handler.NewConsoleHandler(slog.AllLevels)
	h.SetFormatter(f)

	l.AddHandler(h)

	l.WithFields(slog.M{
		"field1": 123,
		"field2": "abc",
	}).Info("info message")

	// PrettyPrint=true

	l = slog.New()
	h = handler.NewConsoleHandler(slog.AllLevels)
	f = slog.NewJSONFormatter(slog.StringMap{
		"level": "levelName",
	})
	f.PrettyPrint = true

	h.SetFormatter(f)
	l.AddHandler(h)

	l.Info("info message and PrettyPrint is TRUE")
}
