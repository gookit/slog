package handler_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func ExampleFileHandler() {
	withLevels := handler.WithLogLevels(slog.Levels{slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel})
	h1 := handler.MustFileHandler("/tmp/error.log", withLevels)

	withLevels = handler.WithLogLevels(slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel})
	h2 := handler.MustFileHandler("/tmp/info.log", withLevels)

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}

func ExampleRotateFileHandler() {
	h1 := handler.MustRotateFile("/tmp/error.log", handler.EveryHour)

	h2 := handler.MustRotateFile("/tmp/info.log", handler.EveryHour)

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}
