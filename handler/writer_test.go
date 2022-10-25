package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestNewIOWriter(t *testing.T) {
	w := new(bytes.Buffer)
	h := handler.NewIOWriter(w, slog.NormalLevels)

	assert.True(t, h.IsHandling(slog.NoticeLevel))

	r := newLogRecord("test io.writer handler")
	assert.NoErr(t, h.Handle(r))
	assert.NoErr(t, h.Flush())

	str := w.String()
	assert.Contains(t, str, "test io.writer handler")

	assert.NoErr(t, h.Close())
}

func TestNewSyncCloser(t *testing.T) {
	logfile := "./testdata/sync-closer.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	f, err := handler.QuickOpenFile(logfile)
	assert.NoErr(t, err)

	h := handler.NewSyncCloser(f, slog.DangerLevels)

	assert.True(t, h.IsHandling(slog.WarnLevel))
	assert.False(t, h.IsHandling(slog.InfoLevel))

	r := newLogRecord("test sync closer handler")
	r.Level = slog.ErrorLevel

	err = h.Handle(r)
	assert.NoErr(t, err)
	assert.NoErr(t, h.Flush())

	bts := fsutil.MustReadFile(logfile)
	str := string(bts)

	assert.Contains(t, str, "test sync closer handler")

	err = h.Close()
	assert.NoErr(t, err)
}

func TestNewWriteCloser(t *testing.T) {
	w := &closedBuffer{Buffer: bytes.Buffer{}}
	h := handler.NewWriteCloser(w, slog.NormalLevels)

	assert.True(t, h.IsHandling(slog.NoticeLevel))

	r := newLogRecord("test writeCloser handler")
	assert.NoErr(t, h.Handle(r))
	assert.NoErr(t, h.Flush())

	str := w.String()
	assert.Contains(t, str, "test writeCloser handler")

	assert.NoErr(t, h.Close())
}

func TestNewFlushCloser(t *testing.T) {
	w := &closedBuffer{}
	h := handler.NewFlushCloser(w, slog.AllLevels)

	r := newLogRecord("TestNewFlushCloser")
	assert.NoErr(t, h.Handle(r))
	assert.NoErr(t, h.Flush())

	str := w.String()
	assert.Contains(t, str, "TestNewFlushCloser")

	assert.NoErr(t, h.Close())
}
