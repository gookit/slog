package main

import (
	"errors"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// run: go run ./_example/slog_all_level.go
func main() {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
	})

	l.AddHandler(handler.NewConsoleHandler(slog.AllLevels))
	printAllLevel(l, "this is a", "log", "message")
}

func printAllLevel(l *slog.Logger, args ...interface{}) {
	l.Debug(args...)
	l.Info(args...)
	l.Warn(args...)
	l.Error(args...)
	l.Print(args...)
	l.Fatal(args...)
	l.Panic(args...)

	l.Trace(args...)
	l.Notice(args...)
	l.ErrorT(errors.New("a error object"))
	l.ErrorT(errorx.New("error with stack info"))
}
