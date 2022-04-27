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
	testFile := "testdata/file.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewFileHandler(testFile)
	assert.NoError(t, err)

	l := slog.NewWithHandlers(h)
	l.PanicFunc = slog.DoNothingOnPanic
	l.Info("info message")
	l.Warn("warn message")
	logAllLevel(l, "file handler message")

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

func TestNewFileHandler_basic(t *testing.T) {
	testFile := "testdata/file_basic.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewFileHandler(testFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, h.Writer())

	r := newLogRecord("test file handler")

	err = h.Handle(r)
	assert.NoError(t, err)

	err = h.Close()
	assert.NoError(t, err)
}

func TestNewSimpleFileHandler(t *testing.T) {
	logfile := "./testdata/simple-file.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))
	assert.False(t, fsutil.IsFile(logfile))

	h, err := handler.NewSimpleFileHandler(logfile)
	assert.NoError(t, err)

	l := slog.NewWithHandlers(h)
	l.Info("info message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(logfile))
	// assert.NoError(t, os.Remove(logfile))
	bts, err := ioutil.ReadFile(testFile)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, slog.WarnLevel.Name())
}
