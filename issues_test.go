package slog_test

import (
	"testing"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestIssues_27(t *testing.T) {
	defer slog.Reset()

	count := 0
	for {
		if count >= 6 {
			break
		}
		slog.Infof("info log %d", count)
		time.Sleep(time.Second)
		count++
	}
}

// https://github.com/gookit/slog/issues/31
func TestIssues_31(t *testing.T) {
	defer slog.Reset()

	h1 := handler.
		MustFileHandler("testdata/error.log", true).
		Configure(func(h *handler.FileHandler) {
			h.BuffSize = 10
		})
	// slog.DangerLevels equals slog.Levels{slog.PanicLevel, slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}
	h1.Levels = slog.DangerLevels

	h2 := handler.
		MustFileHandler("testdata/info.log", true).
		Configure(func(h *handler.FileHandler) {
			h.BuffSize = 10
		})
	h2.Levels = slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}
