package handler_test

import (
	"os"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestNewBufferedHandler(t *testing.T) {
	logfile := "./testdata/buffer-os-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	file, err := handler.QuickOpenFile(logfile)
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	bh := handler.NewBuffered(file, 128)

	// new logger
	l := slog.NewWithHandlers(bh)
	l.Info("buffered info message")

	bts, err := os.ReadFile(logfile)
	assert.NoErr(t, err)
	assert.Empty(t, bts)

	l.Warn("buffered warn message")
	bts, err = os.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")

	err = l.FlushAll()
	assert.NoErr(t, err)
}

func TestLineBufferedFile(t *testing.T) {
	logfile := "./testdata/line-buff-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h, err := handler.LineBufferedFile(logfile, 12, slog.AllLevels)
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	r := newLogRecord("Test LineBufferedFile")
	err = h.Handle(r)
	assert.NoErr(t, err)

	bts, err := os.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "Test LineBufferedFile")
}

func TestLineBuffOsFile(t *testing.T) {
	logfile := "./testdata/line-buff-os-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	file, err := fsutil.QuickOpenFile(logfile)
	assert.NoErr(t, err)

	h := handler.LineBuffOsFile(file, 12, slog.AllLevels)
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	r := newLogRecord("Test LineBuffOsFile")
	err = h.Handle(r)
	assert.NoErr(t, err)

	bts, err := os.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "Test LineBuffOsFile")

	assert.Panics(t, func() {
		handler.LineBuffOsFile(nil, 12, slog.AllLevels)
	})
}

func TestLineBuffWriter(t *testing.T) {
	logfile := "./testdata/line-buff-writer.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	file, err := fsutil.QuickOpenFile(logfile)
	assert.NoErr(t, err)

	h := handler.LineBuffWriter(file, 12, slog.AllLevels)
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))
	assert.Panics(t, func() {
		handler.LineBuffWriter(nil, 12, slog.AllLevels)
	})

	r := newLogRecord("Test LineBuffWriter")
	err = h.Handle(r)
	assert.NoErr(t, err)

	bts, err := os.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "Test LineBuffWriter")

	assert.Panics(t, func() {
		handler.LineBuffOsFile(nil, 12, slog.AllLevels)
	})
}
