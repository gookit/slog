package handler_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func ExampleNewMultiFileHandler() {
	h := &handler.MultiFileHandler{
		FileDir: "testdata/multifiles",
		FileLevels: map[string]slog.Levels{
			"error.log": {slog.ErrorLevel, slog.WarnLevel},
			"info.log":  {slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel},
		},
	}

	slog.AddHandler(h)

	// add logs
	slog.Info("info messages")
}
