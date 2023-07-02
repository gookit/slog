package slog_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// https://github.com/gookit/slog/issues/27
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
	defer slog.MustFlush()

	// slog.DangerLevels equals slog.Levels{slog.PanicLevel, slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}
	h1 := handler.MustFileHandler("testdata/error_issue31.log", handler.WithLogLevels(slog.DangerLevels))

	infoLevels := slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
	h2 := handler.MustFileHandler("testdata/info_issue31.log", handler.WithLogLevels(infoLevels))

	slog.PushHandler(h1)
	slog.PushHandlers(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}

// https://github.com/gookit/slog/issues/52
func TestIssues_52(t *testing.T) {
	testTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}"
	slog.SetLogLevel(slog.ErrorLevel)
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(testTemplate)

	slog.Error("Error message")
	slog.Reset()

	fmt.Println()
	// dump.P(slog.GetFormatter())
}

// https://github.com/gookit/slog/issues/75
func TestIssues_75(t *testing.T) {
	slog.Error("Error message 1")

	// set max level
	slog.SetLogLevel(slog.Level(0))
	// slog.SetLogLevel(slog.PanicLevel)
	slog.Error("Error message 2")
	slog.Reset()
	// dump.P(slog.GetFormatter())
}

// https://github.com/gookit/slog/issues/105
func TestIssues_105(t *testing.T) {
	t.Run("simple write", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			slog.Error("simple error log", i)
			time.Sleep(time.Millisecond * 100)
		}
	})

	// test concurrent write
	t.Run("concurrent write", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				slog.Error("concurrent error log", i)
				time.Sleep(time.Millisecond * 100)
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
