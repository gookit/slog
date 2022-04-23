package handler_test

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestConsoleHandlerWithColor(t *testing.T) {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	l.Configure(func(l *slog.Logger) {
		l.ReportCaller = true
		l.ExitFunc = slog.DoNothingOnExit
	})

	l.Trace("this is a simple log message")
	l.Debug("this is a simple log message")
	l.Info("this is a simple log message")
	l.Notice("this is a simple log message")
	l.Warn("this is a simple log message")
	l.Error("this is a simple log message")
	l.Fatal("this is a simple log message")
}

func TestConsoleHandlerNoColor(t *testing.T) {
	h := handler.NewConsole(slog.AllLevels)
	// no color
	h.TextFormatter().EnableColor = false

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	l.Trace("this is a simple log message")
	l.Debug("this is a simple log message")
	l.Info("this is a simple log message")
	l.Notice("this is a simple log message")
	l.Warn("this is a simple log message")
	l.Error("this is a simple log message")
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

	l.FlushAll()
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

	l.FlushAll()
}
