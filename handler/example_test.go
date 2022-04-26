package handler_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func Example_fileHandler() {
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

func Example_rotateFileHandler() {
	h1 := handler.MustRotateFile("/tmp/error.log", handler.EveryHour, handler.WithLogLevels(slog.DangerLevels))
	h2 := handler.MustRotateFile("/tmp/info.log", handler.EveryHour, handler.WithLogLevels(slog.NormalLevels))

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}
