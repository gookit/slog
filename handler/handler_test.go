package handler_test

import (
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

func TestConsoleHandlerWithColor(t *testing.T) {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	l.DoNothingOnPanicFatal()
	l.Configure(func(l *slog.Logger) {
		l.ReportCaller = true
	})

	logAllLevel(l, "this is a simple log message")
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
	fpath := "./testdata/buffered-os-file.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(fpath))

	file, err := handler.QuickOpenFile(fpath)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	bh := handler.NewBuffered(file, 128)

	// new logger
	l := slog.NewWithHandlers(bh)
	l.Info("buffered info message")

	bts, err := ioutil.ReadFile(fpath)
	assert.NoError(t, err)
	assert.Empty(t, bts)

	l.Warn("buffered warn message")
	bts, err = ioutil.ReadFile(fpath)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")

	err = l.FlushAll()
	assert.NoError(t, err)
}

func TestBufferWrapper(t *testing.T) {
	fpath := "./testdata/buffer-wrapper-simple-file.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(fpath))

	h, err := handler.NewSimpleFile(fpath)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	bw := handler.BufferWrapper(h, 128)

	// new logger
	l := slog.NewWithHandlers(bw)
	l.Info("buffered info message")

	bts, err := ioutil.ReadFile(fpath)
	assert.NoError(t, err)
	assert.Empty(t, bts)

	l.Warn("buffered warn message")
	bts, err = ioutil.ReadFile(fpath)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")

	err = l.FlushAll()
	assert.NoError(t, err)
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
		Extra: map[string]interface{}{
			"source":     "linux",
			"extra_key0": "hello",
			"sub":        slog.M{"sub_key1": "val0"},
		},
	}

	r.Init(false)
	return r
}
