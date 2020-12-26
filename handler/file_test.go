package handler_test

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

const testFile = "./testdata/app.log"
const testSubFile = "./testdata/subdir/app.log"

func TestNewFileHandler(t *testing.T) {
	assert.NoError(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewFileHandler(testFile, false)
	assert.NoError(t, err)

	h.Configure(func(h *handler.FileHandler) {
		h.NoBuffer = true
	})

	l := slog.NewWithHandlers(h)
	l.Info("info message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(testFile))

	bts, err := ioutil.ReadFile(testFile)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARNING]")
	assert.Contains(t, str, "warn message")

	// assert.NoError(t, os.Remove(testFile))
}

func TestNewSimpleFileHandler(t *testing.T) {
	fpath := "./testdata/simple-file.log"
	h, err := handler.NewSimpleFileHandler(fpath)
	assert.NoError(t, err)

	l := slog.NewWithHandlers(h)
	l.Info("info message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(fpath))
	// assert.NoError(t, os.Remove(fpath))
}

func TestNewTimeRotateFileHandler(t *testing.T) {
	fpath := "./testdata/time-rotate-file.log"
	h, err := handler.NewTimeRotateFileHandler(fpath, handler.EveryMinute)
	assert.NoError(t, err)

	l := slog.NewWithHandlers(h)
	l.Info("info", "message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(fpath))
	// assert.NoError(t, os.Remove(fpath))
}