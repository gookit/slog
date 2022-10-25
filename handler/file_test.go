package handler_test

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// const testSubFile = "./testdata/subdir/app.log"

func TestNewFileHandler(t *testing.T) {
	testFile := "testdata/file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewFileHandler(testFile)
	assert.NoErr(t, err)

	l := slog.NewWithHandlers(h)
	l.DoNothingOnPanicFatal()
	l.Info("info message")
	l.Warn("warn message")
	logAllLevel(l, "file handler message")

	assert.True(t, fsutil.IsFile(testFile))

	bts, err := ioutil.ReadFile(testFile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARNING]")
	assert.Contains(t, str, "warn message")

	// assert.NoErr(t, os.Remove(testFile))
}

func TestMustFileHandler(t *testing.T) {
	testFile := "testdata/file-must.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	h := handler.MustFileHandler(testFile)
	assert.NotEmpty(t, h.Writer())

	r := newLogRecord("test file must handler")

	err := h.Handle(r)
	assert.NoErr(t, err)
	assert.NoErr(t, h.Close())

	bts := fsutil.MustReadFile(testFile)
	str := string(bts)

	assert.Contains(t, str, `INFO`)
	assert.Contains(t, str, `test file must handler`)
}

func TestNewFileHandler_basic(t *testing.T) {
	testFile := "testdata/file-basic.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewFileHandler(testFile)
	assert.NoErr(t, err)
	assert.NotEmpty(t, h.Writer())

	r := newLogRecord("test file handler")

	err = h.Handle(r)
	assert.NoErr(t, err)
	assert.NoErr(t, h.Close())

	bts := fsutil.MustReadFile(testFile)
	str := string(bts)

	assert.Contains(t, str, `INFO`)
	assert.Contains(t, str, `test file handler`)
}

func TestNewBuffFileHandler(t *testing.T) {
	testFile := "testdata/file-buff.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.NewBuffFileHandler(testFile, 56)
	assert.NoErr(t, err)
	assert.NotEmpty(t, h.Writer())

	r := newLogRecord("test file buff handler")

	err = h.Handle(r)
	assert.NoErr(t, err)
	assert.NoErr(t, h.Close())

	bts := fsutil.MustReadFile(testFile)
	str := string(bts)

	assert.Contains(t, str, `INFO`)
	assert.Contains(t, str, `test file buff handler`)
}

func TestJSONFileHandler(t *testing.T) {
	testFile := "testdata/file-json.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	h, err := handler.JSONFileHandler(testFile)
	assert.NoErr(t, err)

	r := newLogRecord("test json file handler")
	err = h.Handle(r)
	assert.NoErr(t, err)

	err = h.Close()
	assert.NoErr(t, err)

	bts := fsutil.MustReadFile(testFile)
	str := string(bts)

	assert.Contains(t, str, `"level":"INFO"`)
	assert.Contains(t, str, `"message":"test json file handler"`)
}

func TestMustSimpleFile(t *testing.T) {
	logfile := "./testdata/must-simple-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h := handler.MustSimpleFile(logfile)
	assert.True(t, h.IsHandling(slog.InfoLevel))
}

func TestNewSimpleFileHandler(t *testing.T) {
	logfile := "./testdata/simple-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))
	assert.False(t, fsutil.IsFile(logfile))

	h, err := handler.NewSimpleFileHandler(logfile)
	assert.NoErr(t, err)

	l := slog.NewWithHandlers(h)
	l.Info("info message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(logfile))
	// assert.NoErr(t, os.Remove(logfile))
	bts, err := ioutil.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, slog.WarnLevel.Name())
}
