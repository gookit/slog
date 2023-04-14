package slog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestRecord_AddData(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
		l.CallerFlag = slog.CallerFlagFull
	})
	l.SetHandlers([]slog.Handler{
		handler.NewIOWriter(w, slog.AllLevels),
	})

	r := l.Record()

	// - add data
	r.AddData(testData1).Trace("log message with data")
	s := w.StringReset()
	fmt.Print(s)

	assert.Contains(t, s, "slog_test.TestRecord_AddData")
	assert.Contains(t, s, "log message with data")

	r.AddData(slog.M{"key01": "val01"}).Print("log message add data2")
	s = w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "log message add data2")
	assert.Contains(t, s, "key01:val01")

	// - add value
	r.AddValue("key01", "val02").Println("log message add value")
	s = w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "log message add value")
	assert.Contains(t, s, "key01:val02")

	// -with data
	r.CallerFlag = slog.CallerFlagFcName
	r.WithData(slog.M{"key1": "val1"}).Warn("warn message with data")
	s = w.StringReset()
	fmt.Print(s)

	assert.Contains(t, s, "TestRecord_AddData")
	assert.Contains(t, s, "warn message with data")
	assert.Contains(t, s, "{key1:val1}")
}

func TestRecord_AddExtra(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
		l.CallerFlag = slog.CallerFlagFcName
	})
	l.SetHandlers([]slog.Handler{
		handler.NewIOWriter(w, slog.AllLevels),
	})

	r := l.Record()

	r.AddExtra(testData1).Trace("log message and add extra")
	s := w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "TestRecord_AddExtra")
	assert.Contains(t, s, "log message and add extra")
	assert.Contains(t, s, "key0:val0")

	r.AddExtra(slog.M{"key002": "val002"}).Trace("log message and add extra2")
	s = w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "log message and add extra2")
	assert.Contains(t, s, "TestRecord_AddExtra")
	assert.Contains(t, s, "key002:val002")
}

func TestRecord_SetContext(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
	}).Config(func(l *slog.Logger) {
		l.CallerFlag = slog.CallerFlagPkg
	})
	l.SetHandlers([]slog.Handler{
		handler.NewIOWriter(w, slog.AllLevels),
	})

	r := l.Record()
	r.SetCtx(context.Background()).Info("info message")
	r.WithCtx(context.Background()).Debug("debug message")

	s := w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "github.com/gookit/slog_test")
}

func TestRecord_WithError(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.CallerFlag = slog.CallerFlagFunc
		l.DoNothingOnPanicFatal()
	})
	h := handler.NewIOWriter(w, slog.AllLevels)
	h.SetFormatter(slog.NewTextFormatter("ts={{timestamp}} err={{error}} msg={{message}}\n"))
	l.SetHandlers([]slog.Handler{h})

	r := l.Record()
	r.WithError(errorx.Raw("error message")).Notice("test record with error")

	s := w.StringReset()
	assert.Contains(t, s, "err=error message")
	assert.Contains(t, s, "msg=test record with error")
	fmt.Print(s)
}

func TestRecord_WithTime(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.CallerFlag = slog.CallerFlagFnLine
		l.DoNothingOnPanicFatal()
	})
	l.SetHandlers([]slog.Handler{
		handler.NewIOWriter(w, slog.AllLevels),
	})

	ht := timex.NowHourStart()

	r := l.Record()
	r.WithTime(ht).Notice("a message with time")
	s := w.StringReset()

	assert.Contains(t, s, "a message with time")
	assert.Contains(t, s, timex.FormatByTpl(ht, slog.DefaultTimeFormat))
	fmt.Print(s)
}

func TestRecord_AddFields(t *testing.T) {
	r := newLogRecord("AddFields")

	r.AddFields(slog.M{"f1": "hi", "env": "prod"})
	assert.NotEmpty(t, r.Fields)

	r.AddFields(slog.M{"app": "goods"})
	assert.NotEmpty(t, r.Fields)
}

func TestRecord_SetFields(t *testing.T) {
	r := newLogRecord("AddFields")

	r.SetTime(timex.Now().Yesterday().T())
	r.SetFields(slog.M{"f1": "hi", "env": "prod"})
	assert.NotEmpty(t, r.Fields)
	assert.NotEmpty(t, r.Time)
}

func TestRecord_allLevel(t *testing.T) {
	w := newBuffer()
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
	})
	l.SetHandlers([]slog.Handler{
		handler.NewIOWriter(w, slog.AllLevels),
	})

	r := l.Record()
	r.WithContext(context.Background())
	printAllLevelLogs(r, "a message use record.XX()")
	r.Log(slog.InfoLevel, "a message use record.XX()")
	r.Notice("a message use record.XX()")
	r.Trace("a message use record.XX()")

	s := w.StringReset()
	assert.Contains(t, s, "printAllLevelLogs")
	assert.Contains(t, s, "a message use record.XX()")
	assert.Contains(t, s, "[NOTICE]")
	assert.Contains(t, s, "[TRACE]")

	printfAllLevelLogs(r, "a message use %s()", "record.XXf")
	r.Logf(slog.InfoLevel, "a message use %s()", "record.XXf")
	r.Noticef("a message use %s()", "record.XXf")
	r.Tracef("a message use %s()", "record.XXf")

	s = w.StringReset()
	assert.Contains(t, s, "printfAllLevelLogs")
	assert.Contains(t, s, "a message use record.XXf()")
	assert.Contains(t, s, "[NOTICE]")
	assert.Contains(t, s, "[TRACE]")
}
