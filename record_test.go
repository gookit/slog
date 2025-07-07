package slog_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gookit/goutil/byteutil"
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
	assert.Eq(t, "val01", r.Value("key01"))

	// - add value
	r.AddValue("key01", "val02").Println("log message add value")
	s = w.StringReset()
	fmt.Print(s)
	assert.Contains(t, s, "log message add value")
	assert.Contains(t, s, "key01:val02")
	// - first add value
	nr := &slog.Record{}
	assert.Nil(t, nr.Value("key01"))
	nr.WithValue("key01", "val02")
	assert.Eq(t, "val02", nr.Value("key01"))

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

	r.AddExtra(slog.M{"key002": "val002"}).AddExtra(slog.M{"key01": "val01"}).
		Trace("log message and add extra2")
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
	w, l := newTestLogger()
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

	// WithFields
	r = r.WithFields(slog.M{"f2": "v2"})
	assert.Eq(t, "v2", r.Field("f2"))

	// - first add field
	nr := slog.Record{}
	assert.Nil(t, nr.Field("f3"))
	nr.AddField("f3", "val02")
	assert.Eq(t, "val02", nr.Field("f3"))
}

func TestRecord_WithFields(t *testing.T) {
	w, l := newTestLogger()
	r := l.Record().
		WithFields(slog.M{"key1": "value1", "key2": "value2"}).
		WithFields(slog.M{"key3": "value3"})
	assert.Eq(t, "value1", r.Field("key1"))
	assert.Eq(t, "value2", r.Field("key2"))
	assert.Eq(t, "value3", r.Field("key3"))

	r.Info("log message with fields")
	s := w.StringReset()
	fmt.Print(s)

	assert.Contains(t, s, "log message with fields")
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
	r = r.WithContext(context.Background())
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

func TestRecord_useMultiTimes(t *testing.T) {
	buf := byteutil.NewBuffer()
	l := slog.NewWithHandlers(
		handler.NewSimple(buf, slog.DebugLevel),
		handler.NewSimple(os.Stdout, slog.DebugLevel),
	)

	r := l.Record()
	t.Run("simple", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			r.Error("simple error log", i)
			time.Sleep(time.Millisecond * 100)
		}
	})

	// test concurrent write
	t.Run("concurrent", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				r.Error("concurrent error log", i)
				time.Sleep(time.Millisecond * 100)
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
