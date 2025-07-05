package slog_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var doNothing = func(code int) {
	// do nothing
}

func TestStd(t *testing.T) {
	defer slog.Reset()
	assert.Eq(t, "stdLogger", slog.Std().Name())

	_, ok := slog.GetFormatter().(*slog.TextFormatter)
	assert.True(t, ok)

	slog.SetLogLevel(slog.WarnLevel)
	slog.SetFormatter(slog.NewJSONFormatter())

	assert.True(t, slog.Std().IsHandling(slog.WarnLevel))
	assert.True(t, slog.Std().IsHandling(slog.ErrorLevel))
	assert.False(t, slog.Std().IsHandling(slog.InfoLevel))

	_, ok = slog.GetFormatter().(*slog.JSONFormatter)
	assert.True(t, ok)

	buf := new(bytes.Buffer)
	slog.Std().ExitFunc = func(code int) {
		buf.WriteString("Exited,")
		buf.WriteString(strconv.Itoa(code))
	}
	slog.Exit(34)
	assert.Eq(t, "Exited,34", buf.String())
}

func TestTextFormatNoColor(t *testing.T) {
	defer slog.Reset()
	slog.Configure(func(l *slog.SugaredLogger) {
		f := l.Formatter.(*slog.TextFormatter)
		f.EnableColor = false

		l.DoNothingOnPanicFatal()
	})

	printLogs("print log message")
	printfLogs("print log with %s", "params")

	assert.NoErr(t, slog.Std().FlushAll())
	assert.NoErr(t, slog.Std().Close())
}

func TestFlushDaemon(t *testing.T) {
	defer slog.Reset()

	buf := byteutil.NewBuffer()
	slog.Configure(func(l *slog.SugaredLogger) {
		l.FlushInterval = timex.Millisecond * 100
		l.Output = buf
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	go slog.FlushDaemon(func() {
		fmt.Println("flush daemon stopped")
		wg.Done()
	})

	go func() {
		// mock app running
		time.Sleep(time.Second * 1)

		// stop daemon
		fmt.Println("stop flush daemon")
		slog.StopDaemon()
	}()

	slog.Info("print log message")

	wg.Wait()

	fmt.Print(buf.ResetGet())
}

func TestFlushTimeout(t *testing.T) {
	defer slog.Reset()
	slog.Info("print log message")
	slog.FlushTimeout(timex.Second * 1)
	slog.MustFlush()
}

func TestNewSugaredLogger(t *testing.T) {
	buf := byteutil.NewBuffer()
	l := slog.NewSugared(buf, slog.DebugLevel, func(sl *slog.SugaredLogger) {
		sl.SetName("test")
		sl.ReportCaller = true
		sl.CallerFlag = slog.CallerFlagFcLine
	})

	l.Debug("debug message")
	l.Info("info message")
	s := buf.ResetAndGet()
	assert.StrContains(t, s, "debug message")

	l = slog.NewStd(func(sl *slog.SugaredLogger) {
		sl.SetName("test")
		sl.ReportCaller = true
		sl.CallerFlag = slog.CallerFlagFunc
	})
	l.Info("info message1")
}

type logTest struct {
	*slog.SugaredLogger
}

func (l logTest) testPrint() {
	l.Logger.Info("print testing")
}

func TestTextFormatWithColor(t *testing.T) {
	defer slog.Reset()

	slog.Configure(func(l *slog.SugaredLogger) {
		l.Level = slog.TraceLevel
		l.DoNothingOnPanicFatal()
	})

	printLogs("this is a simple log message")
	fmt.Println()

	slog.Std().Trace("this is a simple log message")
	lt := &logTest{slog.Std()}
	lt.testPrint()

	fmt.Println()
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(slog.NamedTemplate)
	printfLogs("print log with %s", "params")

	fmt.Println()
	tpl := "[{{datetime}}] [{{channel}}] [{{level}}] [{{func}}] {{message}} {{data}} {{extra}}\n"
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(tpl)
	printfLogs("print log with %s", "params")

	lt = &logTest{
		slog.Std(),
	}
	lt.testPrint()
}

func printLogs(msg string) {
	slog.Log(slog.TraceLevel, msg)
	slog.Print(msg)
	slog.Println(msg)
	slog.Trace(msg)
	slog.Debug(msg)
	slog.Info(msg)
	slog.Notice(msg)
	slog.Warn(msg)
	slog.Error(msg)
	slog.Fatal(msg)
	slog.FatalErr(errorx.Rawf("Fatal Err: %s", msg))
	slog.Panic(msg)
	slog.PanicErr(errorx.Rawf("Panic Err: %s", msg))
	slog.ErrorT(errors.New(msg))
	slog.ErrorT(errorx.Newf("Traced Err: %s", msg))
}

func printfLogs(msg string, args ...any) {
	slog.Printf(msg, args...)
	slog.Tracef(msg, args...)
	slog.Debugf(msg, args...)
	slog.Infof(msg, args...)
	slog.Noticef(msg, args...)
	slog.Warnf(msg, args...)
	slog.Errorf(msg, args...)
	slog.Panicf(msg, args...)
	slog.Fatalf(msg, args...)
}

func TestSetFormatter_jsonFormat(t *testing.T) {
	defer slog.Reset()
	slog.SetLevelByName("trace")
	slog.SetFormatter(slog.NewJSONFormatter())

	th := newTestHandler()
	th.SetFormatter(slog.NewJSONFormatter().Configure(func(f *slog.JSONFormatter) {
		f.Fields = slog.NoTimeFields
	}))
	slog.PushHandler(th)

	assert.Eq(t, 2, slog.Std().HandlersNum())

	slog.Info("info log message1")
	slog.Warn("warning log message2")
	s := th.ResetGet()
	assert.StrContains(t, s, `"level":"INFO"`)
	assert.StrContains(t, s, `info log message1`)
	assert.StrContains(t, s, `"level":"WARNING"`)
	assert.StrContains(t, s, `warning log message2`)

	slog.WithData(slog.M{
		"key0": 134,
		"key1": "abc",
	}).Infof("info log %s", "message")
	s = th.ResetGet()

	r := slog.WithFields(slog.M{
		"category": "service",
		"IP":       "127.0.0.1",
	})
	r.Infof("info %s", "message")
	r.Debugf("debug %s", "message")
	s = th.ResetGet()

	r = slog.WithField("app", "order")
	r.Trace("trace message")
	r.Println("print message")
	s = th.ResetGet()
	assert.StrContains(t, s, `"app":"order"`)
	assert.StrCount(t, s, `"app":"order"`, 2)

	slog.WithContext(context.Background()).Print("print message with ctx")
	assert.StrContains(t, th.ResetGet(), "print message with ctx")
}

func TestAddHandler(t *testing.T) {
	defer slog.Reset()
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(slog.NewJSONFormatter().Configure(func(f *slog.JSONFormatter) {
		f.Aliases = slog.StringMap{
			"level":   "levelName",
			"message": "msg",
			"data":    "params",
		}
	}))

	slog.AddHandlers(h2)
	slog.Infof("info %s", "message")
}

func TestWithExtra(t *testing.T) {
	defer slog.Reset()

	th := newTestHandler()
	slog.AddHandler(th)

	slog.WithExtra(slog.M{"ext1": "val1"}).
		AddValue("key1", "val2").
		Info("info message")
	s := th.ResetGet()
	assert.StrContains(t, s, `ext1:val1`)
	assert.StrContains(t, s, `{key1:val2}`)

	slog.WithValue("key1", "val2").Info("info message")
	s = th.ResetGet()
	assert.StrContains(t, s, `{key1:val2}`)
}

func TestAddProcessor(t *testing.T) {
	defer slog.Reset()

	buf := new(bytes.Buffer)
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Level = slog.TraceLevel
		logger.Output = buf
		logger.Formatter = slog.NewJSONFormatter()
	})

	slog.AddProcessor(slog.AddHostname())
	slog.Trace("Trace message")
	slog.Tracef("Tracef %s", "message")

	str := buf.String()
	buf.Reset()
	fmt.Println(str)
	assert.Contains(t, str, `"hostname":`)
	assert.Contains(t, str, "Trace message")
	assert.Contains(t, str, "Tracef message")

	slog.AddProcessors(slog.ProcessorFunc(func(r *slog.Record) {
		r.AddField("newField", "newValue")
	}))
	slog.Debug("Debug message")
	slog.Debugf("Debugf %s", "message")
	str = buf.String()
	buf.Reset()

	assert.Contains(t, str, `"newField":"newValue"`)
	assert.Contains(t, str, "Debug message")
	assert.Contains(t, str, "Debugf message")
}


func TestPrependExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	assert.Len(t, slog.ExitHandlers(), 2)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Eq(t, "HANDLER2-HANDLER1-Exited", buf.String())
}

func TestRegisterExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	// prepend
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER3-")
	})
	assert.Len(t, slog.ExitHandlers(), 3)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Eq(t, "HANDLER3-HANDLER1-HANDLER2-Exited", buf.String())
}

func TestExitHandlerWithError(t *testing.T) {
	defer slog.Reset()
	assert.Len(t, slog.ExitHandlers(), 0)

	slog.RegisterExitHandler(func() {
		panic("test error")
	})

	slog.SetExitFunc(func(code int) {})

	testutil.RewriteStderr()
	slog.Exit(23)
	str := testutil.RestoreStderr()
	assert.Eq(t, "slog: run exit handler(global) recovered, error: test error\n", str)
}

func TestLogger_ExitHandlerWithError(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ExitFunc = doNothing
	})

	assert.Len(t, l.ExitHandlers(), 0)

	l.RegisterExitHandler(func() {
		panic("test error")
	})

	testutil.RewriteStderr()
	l.Exit(23)
	str := testutil.RestoreStderr()
	assert.Eq(t, "slog: run exit handler recovered, error: test error\n", str)
}

func TestLogger_PrependExitHandler(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ExitFunc = doNothing
	})

	assert.Len(t, l.ExitHandlers(), 0)

	l.PrependExitHandler(func() {
		panic("test error2")
	})

	testutil.RewriteStderr()
	l.Exit(23)
	str := testutil.RestoreStderr()
	assert.Eq(t, "slog: run exit handler recovered, error: test error2\n", str)
}

func TestSugaredLogger_Close(t *testing.T) {
	h := newTestHandler()

	sl := slog.NewStd(func(sl *slog.SugaredLogger) {
		sl.PushHandler(h)
		sl.Formatter = newTestFormatter()
	})

	h.errOnClose = true
	err := sl.Close()
	assert.Err(t, err)
	assert.Err(t, sl.LastErr())
	assert.Eq(t, "close error", err.Error())
}

func TestSugaredLogger_Handle(t *testing.T) {
	buf := byteutil.NewBuffer()
	sl := slog.NewStd(func(sl *slog.SugaredLogger) {
		sl.Output = buf
		sl.Formatter = newTestFormatter(true)
	})

	// Handle error: format error
	sl.WithField("key", "value").Error("error message")
	err := sl.LastErr()
	assert.Err(t, err)
	assert.Eq(t, "format error", err.Error())
}
