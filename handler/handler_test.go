package handler_test

import (
	"log/syslog"
	"testing"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var (
	sampleData = slog.M{
		"name":  "inhere",
		"age":   100,
		"skill": "go,php,java",
	}
)

func TestMain(m *testing.M) {
	goutil.PanicErr(fsutil.RemoveSub("./testdata", fsutil.ExcludeNames(".keep")))
	m.Run()
}

func TestConfig_CreateWriter(t *testing.T) {
	cfg := handler.NewEmptyConfig()

	w, err := cfg.CreateWriter()
	assert.Nil(t, w)
	assert.Err(t, err)

	h, err := cfg.CreateHandler()
	assert.Nil(t, h)
	assert.Err(t, err)

	logfile := "./testdata/file-by-config.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	cfg.With(
		handler.WithBuffMode(handler.BuffModeBite),
		handler.WithLogLevels(slog.NormalLevels),
		handler.WithLogfile(logfile),
	)

	w, err = cfg.CreateWriter()
	assert.NoErr(t, err)

	_, err = w.Write([]byte("hello, config"))
	assert.NoErr(t, err)

	bts := fsutil.MustReadFile(logfile)
	str := string(bts)

	assert.Eq(t, str, "hello, config")
	assert.NoErr(t, w.Sync())
	assert.NoErr(t, w.Close())
}

func TestConfig_RotateWriter(t *testing.T) {
	cfg := handler.NewEmptyConfig()

	w, err := cfg.RotateWriter()
	assert.Nil(t, w)
	assert.Err(t, err)
}

func TestConsoleHandlerWithColor(t *testing.T) {
	l := slog.NewWithHandlers(handler.ConsoleWithLevels(slog.AllLevels))
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

func TestNewEmailHandler(t *testing.T) {
	from := handler.EmailOption{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: 587,
		FromAddr: "someone@gmail.com",
	}

	h := handler.NewEmailHandler(from, []string{
		"another@gmail.com",
	})

	assert.Eq(t, slog.InfoLevel, h.Level)

	// handle error
	h.SetFormatter(newTestFormatter(true))
	assert.Err(t, h.Handle(newLogRecord("test email handler")))
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

func TestNopFlushClose_Flush(t *testing.T) {
	nfc := handler.NopFlushClose{}

	assert.NoErr(t, nfc.Flush())
	assert.NoErr(t, nfc.Close())
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
	assert.Eq(t, 2, a)
}

func TestNewSysLogHandler(t *testing.T) {
	if sysutil.IsWin() {
		t.Skip("skip test on windows")
	}

	h, err := handler.NewSysLogHandler(syslog.LOG_INFO, "slog")
	assert.NoErr(t, err)

	err = h.Handle(newLogRecord("test syslog handler"))
	assert.NoErr(t, err)

	assert.NoErr(t, h.Flush())
	assert.NoErr(t, h.Close())
}

func logAllLevel(log slog.SLogger, msg string) {
	for _, level := range slog.AllLevels {
		log.Log(level, msg)
	}
}

func newLogRecord(msg string) *slog.Record {
	r := &slog.Record{
		Channel: "handler_test",
		Level:   slog.InfoLevel,
		Message: msg,
		Time:    slog.DefaultClockFn.Now(),
		Data:    sampleData,
		Extra: map[string]any{
			"source":     "linux",
			"extra_key0": "hello",
			"sub":        slog.M{"sub_key1": "val0"},
		},
	}

	r.Init(false)
	return r
}

type testHandler struct {
	errOnHandle bool
	errOnFlush  bool
	errOnClose  bool
}

func newTestHandler() *testHandler {
	return &testHandler{}
}

// func (h testHandler) Reset() {
// 	h.errOnHandle = false
// 	h.errOnFlush = false
// 	h.errOnClose = false
// }

func (h testHandler) IsHandling(_ slog.Level) bool {
	return true
}

func (h testHandler) Close() error {
	if h.errOnClose {
		return errorx.Raw("close error")
	}
	return nil
}

func (h testHandler) Flush() error {
	if h.errOnFlush {
		return errorx.Raw("flush error")
	}
	return nil
}

func (h testHandler) Handle(_ *slog.Record) error {
	if h.errOnHandle {
		return errorx.Raw("handle error")
	}
	return nil
}

type testFormatter struct {
	errOnFormat bool
}

func newTestFormatter(errOnFormat ...bool) *testFormatter {
	return &testFormatter{
		errOnFormat: len(errOnFormat) > 0 && errOnFormat[0],
	}
}

func (f testFormatter) Format(r *slog.Record) ([]byte, error) {
	if f.errOnFormat {
		return nil, errorx.Raw("format error")
	}
	return []byte(r.Message), nil
}
