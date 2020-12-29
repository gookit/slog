package handler_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func ExampleFileHandler() {
	h1 := handler.MustFileHandler("/tmp/error.log", true)
	h1.Levels = slog.Levels{slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}

	h2 := handler.MustFileHandler("/tmp/info.log", true)
	h1.Levels = slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}
