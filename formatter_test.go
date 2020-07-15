package slog_test

import (
	"fmt"
	"testing"

	"github.com/gookit/slog/formatter"
	"github.com/gookit/slog/handler"
)

func TestNewLineFormatter(t *testing.T) {
	lf := formatter.NewLineFormatter()

	fmt.Println(lf.FieldMap())
}

func TestJSONFormatter(t *testing.T) {
	l := New()

	f := formatter.NewJSONFormatter(nil)
	h := handler.NewConsoleHandler(AllLevels)
	h.SetFormatter(f)

	l.AddHandler(h)

	l.Info("info message")

	// PrettyPrint=true

	l = New()
	h = handler.NewConsoleHandler(AllLevels)
	f = formatter.NewJSONFormatter(StringMap{
		"level": "levelName",
	})
	f.PrettyPrint = true

	h.SetFormatter(f)
	l.AddHandler(h)

	l.Info("info message and PrettyPrint is TRUE")
}
