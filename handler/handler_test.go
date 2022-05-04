package handler_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

var (
	sampleData = slog.M{
		"name":  "inhere",
		"age":   100,
		"skill": "go,php,java",
	}
)

func TestConfig_CreateWriter(t *testing.T) {
	cfg := handler.NewEmptyConfig()

	w, err := cfg.CreateWriter()
	assert.Nil(t, w)
	assert.Error(t, err)

	h, err := cfg.CreateHandler()
	assert.Nil(t, h)
	assert.Error(t, err)

	logfile := "./testdata/file-by-config.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	cfg.With(
		handler.WithBuffMode(handler.BuffModeBite),
		handler.WithLogLevels(slog.NormalLevels),
		handler.WithLogfile(logfile),
	)

	w, err = cfg.CreateWriter()
	assert.NoError(t, err)

	_, err = w.Write([]byte("hello, config"))
	assert.NoError(t, err)

	bts := fsutil.MustReadFile(logfile)
	str := string(bts)

	assert.Equal(t, str, "hello, config")
	assert.NoError(t, w.Sync())
	assert.NoError(t, w.Close())
}

func TestConfig_RotateWriter(t *testing.T) {
	cfg := handler.NewEmptyConfig()

	w, err := cfg.RotateWriter()
	assert.Nil(t, w)
	assert.Error(t, err)
}

func TestConsoleHandlerWithColor(t *testing.T) {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	l.DoNothingOnPanicFatal()
	l.Configure(func(l *slog.Logger) {
		l.ReportCaller = true
	})

	logAllLevel(l, "this is a simple log message")
	// logfAllLevel()
}

func TestConsoleHandlerNoColor(t *testing.T) {
	h := handler.NewConsole(slog.AllLevels)
	// no color
	h.TextFormatter().EnableColor = false

	l := slog.NewWithHandlers(h)
	l.DoNothingOnPanicFatal()
	l.ReportCaller = true

	logAllLevel(l, "this is a simple log message")
}

func TestNewBufferedHandler(t *testing.T) {
	logfile := "./testdata/buffer-os-file.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	file, err := handler.QuickOpenFile(logfile)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	bh := handler.NewBuffered(file, 128)

	// new logger
	l := slog.NewWithHandlers(bh)
	l.Info("buffered info message")

	bts, err := ioutil.ReadFile(logfile)
	assert.NoError(t, err)
	assert.Empty(t, bts)

	l.Warn("buffered warn message")
	bts, err = ioutil.ReadFile(logfile)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")

	err = l.FlushAll()
	assert.NoError(t, err)
}

func TestBufferWrapper(t *testing.T) {
	logfile := "./testdata/buffer-wrap-handler.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	h, err := handler.NewSimpleFile(logfile)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	bw := handler.BufferWrapper(h, 128)

	// new logger
	l := slog.NewWithHandlers(bw)
	l.Info("buffered info message")

	bts, err := ioutil.ReadFile(logfile)
	assert.NoError(t, err)
	assert.Empty(t, bts)

	l.Warn("buffered warn message")
	bts, err = ioutil.ReadFile(logfile)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")

	err = l.FlushAll()
	assert.NoError(t, err)
}

func TestNewEmailHandler(t *testing.T) {
	from := handler.EmailOption{
		SmtpHost: "smtp.gmail.com",
		SmtpPort: 587,
		FromAddr: "someone@gmail.com",
	}

	h := handler.NewEmailHandler(from, []string{
		"another@gmail.com",
	})

	assert.Equal(t, slog.InfoLevel, h.Level)
}

func TestLevelWithFormatter(t *testing.T) {
	lf := handler.LevelWithFormatter{Level: slog.InfoLevel}

	assert.True(t, lf.IsHandling(slog.ErrorLevel))
	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.False(t, lf.IsHandling(slog.DebugLevel))
}

func TestLevelsWithFormatter(t *testing.T) {
	lsf := handler.LevelsWithFormatter{Levels: slog.NormalLevels}

	assert.False(t, lsf.IsHandling(slog.ErrorLevel))
	assert.True(t, lsf.IsHandling(slog.InfoLevel))
	assert.True(t, lsf.IsHandling(slog.DebugLevel))
}

func TestNewSimpleHandler(t *testing.T) {
	buf := new(bytes.Buffer)

	h := handler.NewSimple(buf, slog.InfoLevel)
	r := newLogRecord("test simple handler")
	assert.NoError(t, h.Handle(r))

	s := buf.String()
	buf.Reset()
	fmt.Print(s)
	assert.Contains(t, s, "test simple handler")

	assert.NoError(t, h.Flush())
	assert.NoError(t, h.Close())

	h = handler.NewHandler(buf, slog.InfoLevel)
	r = newLogRecord("test simple handler2")
	assert.NoError(t, h.Handle(r))

	s = buf.String()
	buf.Reset()
	fmt.Print(s)
	assert.Contains(t, s, "test simple handler2")

	assert.NoError(t, h.Flush())
	assert.NoError(t, h.Close())
}

func TestNopFlushClose_Flush(t *testing.T) {
	nfc := handler.NopFlushClose{}

	assert.NoError(t, nfc.Flush())
	assert.NoError(t, nfc.Close())
}

func TestLockWrapper_Lock(t *testing.T) {
	lw := &handler.LockWrapper{}
	assert.True(t, lw.LockEnabled())

	lw.EnableLock(true)
	assert.True(t, lw.LockEnabled())

	a := 1
	lw.Lock()
	a++
	lw.Unlock()
	assert.Equal(t, 2, a)
}

func logAllLevel(log slog.SLogger, msg string) {
	for _, level := range slog.AllLevels {
		log.Log(level, msg)
	}
}

func logfAllLevel(log slog.SLogger, tpl string, args ...interface{}) {
	for _, level := range slog.AllLevels {
		log.Logf(level, tpl, args...)
	}
}

func newLogRecord(msg string) *slog.Record {
	r := &slog.Record{
		Channel: "handler_test",
		Level:   slog.InfoLevel,
		Message: msg,
		Time:    slog.DefaultClockFn.Now(),
		Data:    sampleData,
		Extra: map[string]interface{}{
			"source":     "linux",
			"extra_key0": "hello",
			"sub":        slog.M{"sub_key1": "val0"},
		},
	}

	r.Init(false)
	return r
}
