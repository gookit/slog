package handler_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func ExampleFileHandler() {
	c1 := handler.NewSimpleConfig(func(c *handler.SimpleConfig) {
		c.Levels = slog.Levels{slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}
	})
	h1 := handler.MustFileHandler("/tmp/error.log", c1)

	c2 := handler.NewSimpleConfig(func(c *handler.SimpleConfig) {
		c.Levels = slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
	})
	h2 := handler.MustFileHandler("/tmp/info.log", c2)

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
