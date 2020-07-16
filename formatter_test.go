package slog_test

import (
	"fmt"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestNewTextFormatter(t *testing.T) {
	lf := slog.NewTextFormatter()

	fmt.Println(lf.FieldMap())
}

func TestJSONFormatter(t *testing.T) {
	l := slog.New()

	f := slog.NewJSONFormatter(nil)
	h := handler.NewConsoleHandler(slog.AllLevels)
	h.SetFormatter(f)

	l.AddHandler(h)

	fields := slog.M{
		"field1": 123,
		"field2": "abc",
	}

	l.WithFields(fields).Info("info message")

	// PrettyPrint=true

	l = slog.New()
	h = handler.NewConsoleHandler(slog.AllLevels)
	f = slog.NewJSONFormatter(slog.StringMap{
		"level": "levelName",
	})
	f.PrettyPrint = true

	h.SetFormatter(f)
	l.AddHandler(h)

	l.WithFields(fields).
		SetData(slog.M{"key1": "val1"}).
		SetExtra(slog.M{"ext1": "val1"}).
		Info("info message and PrettyPrint is TRUE")
}
