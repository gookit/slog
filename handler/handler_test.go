package handler_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var doNothing = func(code int) {
	// do nothing
}

func TestConsoleHandlerWithColor(t *testing.T) {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels)).Configure(func(l *slog.Logger) {
		l.ReportCaller = true
		l.ExitFunc = doNothing
	})

	l.Trace("this is a simple log message")
	l.Debug("this is a simple log message")
	l.Info("this is a simple log message")
	l.Notice("this is a simple log message")
	l.Warn("this is a simple log message")
	l.Error("this is a simple log message")
	l.Fatal("this is a simple log message")
}

func TestConsoleHandlerNoColor(t *testing.T) {

}
