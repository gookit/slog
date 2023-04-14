package slog_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/gsr"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestLoggerBasic(t *testing.T) {
	l := slog.New()
	l.SetName("testName")
	assert.Eq(t, "testName", l.Name())

	l = slog.NewWithName("testName")
	assert.Eq(t, "testName", l.Name())
}

func TestLogger_PushHandler(t *testing.T) {
	l := slog.New().Config(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
	})

	w1 := new(bytes.Buffer)
	h1 := handler.NewIOWriterHandler(w1, slog.DangerLevels)
	l.PushHandler(h1)

	w2 := new(bytes.Buffer)
	h2 := handler.NewIOWriterHandler(w2, slog.NormalLevels)
	l.PushHandlers(h2)

	l.Warning(slog.WarnLevel, "message")
	l.Logf(slog.TraceLevel, "%s message", slog.TraceLevel)

	assert.Contains(t, w1.String(), "WARNING message")
	assert.Contains(t, w2.String(), "TRACE message")
	assert.Contains(t, w2.String(), "TestLogger_PushHandler")

	assert.NoErr(t, l.Sync())
	assert.NoErr(t, l.Flush())
	l.MustFlush()

	assert.NoErr(t, l.Close())
	l.Reset()
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

	l.WithField("newKey", "value").Fatalln("a fatal message")
	l.WithTime(timex.NowHourStart()).Panicln("a panic message")
}

func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	h := handler.NewIOWriterHandler(&buf, slog.AllLevels)

	l := newLogger()
	l.AddHandlers(h)

	ctx := context.Background()

	r := l.WithCtx(ctx)
	r.Info("with context")

	str := buf.String()
	assert.Contains(t, str, `with context`)
}

func TestLogger_panic(t *testing.T) {
	h := newTestHandler()
	h.errOnFlush = true

	l := slog.NewWithHandlers(h)

	assert.Panics(t, func() {
		l.MustFlush()
	})

	err := l.LastErr()
	assert.Err(t, err)
	assert.Eq(t, "flush error", err.Error())
}

func TestLogger_error(t *testing.T) {
	h := newTestHandler()
	l := slog.NewWithHandlers(h)

	err := l.VisitAll(func(h slog.Handler) error {
		return errorx.Raw("visit error")
	})
	assert.Err(t, err)
	assert.Eq(t, "visit error", err.Error())

	h.errOnClose = true
	err = l.Close()
	assert.Err(t, err)
	assert.Eq(t, "close error", err.Error())
}

func TestLogger_panicLevel(t *testing.T) {
	w := new(bytes.Buffer)
	l := slog.NewWithHandlers(handler.NewIOWriter(w, slog.AllLevels))

	// assert.PanicsWithValue(t, "slog: panic message", func() {
	assert.Panics(t, func() {
		l.Panicln("panicln message")
	})
	assert.Contains(t, w.String(), "[PANIC]")
	assert.Contains(t, w.String(), "panicln message")

	w.Reset()
	assert.Panics(t, func() {
		l.Panicf("panicf message")
	})
	assert.Contains(t, w.String(), "panicf message")

	w.Reset()
	assert.Panics(t, func() {
		l.Panic("panic message")
	})
	assert.Contains(t, w.String(), "panic message")

	assert.NoErr(t, l.FlushAll())
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

func TestLogger_write_error(t *testing.T) {
	h := newTestHandler()
	h.errOnHandle = true

	l := slog.NewWithHandlers(h)
	l.Info("a message")

	err := l.LastErr()
	assert.Err(t, err)
	assert.Eq(t, "handle error", err.Error())
}

func newLogger() *slog.Logger {
	return slog.NewWithConfig(func(l *slog.Logger) {
		l.ReportCaller = true
		l.DoNothingOnPanicFatal()
	})
}

func printAllLevelLogs(l gsr.Logger, args ...interface{}) {
	l.Debug(args...)
	l.Info(args...)
	l.Warn(args...)
	l.Error(args...)
	l.Print(args...)
	l.Println(args...)
	l.Fatal(args...)
	l.Fatalln(args...)
	l.Panic(args...)
	l.Panicln(args...)

	sl, ok := l.(*slog.Logger)
	if ok {
		sl.Trace(args...)
		sl.Notice(args...)
		sl.ErrorT(errors.New("a error object"))
		sl.ErrorT(errorx.New("error with stack info"))
	}
}

func printfAllLevelLogs(l gsr.Logger, tpl string, args ...interface{}) {
	l.Printf(tpl, args...)
	l.Debugf(tpl, args...)
	l.Infof(tpl, args...)
	l.Warnf(tpl, args...)
	l.Errorf(tpl, args...)
	l.Panicf(tpl, args...)
	l.Fatalf(tpl, args...)

	if sl, ok := l.(*slog.Logger); ok {
		sl.Noticef(tpl, args...)
		sl.Tracef(tpl, args...)
	}
}
