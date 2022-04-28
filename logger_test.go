package slog_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestLoggerBasic(t *testing.T) {
	l := slog.New()
	l.SetName("testName")
	assert.Equal(t, "testName", l.Name())

	l = slog.NewWithName("testName")
	assert.Equal(t, "testName", l.Name())
}

func TestLogger_AddHandlers(t *testing.T) {

}

func TestLogger_ReportCaller(t *testing.T) {
	l := slog.NewWithConfig(func(logger *slog.Logger) {
		logger.ReportCaller = true
		logger.CallerFlag = slog.CallerFlagFnLine
	})

	var buf bytes.Buffer
	h := handler.NewIOWriterHandler(&buf, slog.AllLevels)
	h.SetFormatter(slog.NewJSONFormatter(func(f *slog.JSONFormatter) {
		f.Fields = append(f.Fields, slog.FieldKeyCaller)
	}))

	l.AddHandler(h)
	l.Info("message")

	str := buf.String()
	assert.Contains(t, str, `"caller":"logger_test.go`)
}

func TestLogger_Log(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ReportCaller = true
		l.DoNothingOnPanicFatal()
	})

	l.AddHandler(handler.NewConsoleHandler(slog.AllLevels))
	l.Log(slog.InfoLevel, "a", slog.InfoLevel, "message")
}

func TestLogger_log_allLevel(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ReportCaller = true
		l.DoNothingOnPanicFatal()
	})

	l.AddHandler(handler.NewConsoleHandler(slog.AllLevels))
	printAllLevelLogs(l, "this a", "log", "message")
}

func TestLogger_logf_allLevel(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ReportCaller = true
		l.CallerFlag = slog.CallerFlagFpLine
		l.DoNothingOnPanicFatal()
	})

	l.AddHandler(handler.NewConsoleHandler(slog.AllLevels))
	printfAllLevelLogs(l, "this a log %s", "message")
}

func printAllLevelLogs(l *slog.Logger, args ...interface{}) {
	l.Print(args...)
	l.Println(args...)
	l.Trace(args...)
	l.Debug(args...)
	l.Info(args...)
	l.Notice(args...)
	l.Warn(args...)
	l.Error(args...)
	l.Fatal(args...)
	l.Panic(args...)
	l.ErrorT(errors.New("a error object"))
	l.ErrorT(errorx.New("error with stack info"))
}

func printfAllLevelLogs(l *slog.Logger, tpl string, args ...interface{}) {
	l.Printf(tpl, args...)
	l.Tracef(tpl, args...)
	l.Debugf(tpl, args...)
	l.Infof(tpl, args...)
	l.Noticef(tpl, args...)
	l.Warnf(tpl, args...)
	l.Errorf(tpl, args...)
	l.Panicf(tpl, args...)
	l.Fatalf(tpl, args...)
}
