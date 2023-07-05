package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/testutil/fakeobj"
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
	logfile := "./testdata/sync_closer.log"

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
	w := fakeobj.NewWriter()
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
	w := fakeobj.NewWriter()
	h := handler.NewFlushCloser(w, slog.AllLevels)
	w.WriteString("before flush\n")

	r := newLogRecord("TestNewFlushCloser")
	assert.NoErr(t, h.Handle(r))

	str := w.ResetGet()
	assert.Contains(t, str, "TestNewFlushCloser")

	assert.NoErr(t, h.Flush())
	assert.NoErr(t, h.Close())

	h = handler.FlushCloserWithMaxLevel(w, slog.WarnLevel)
	r = newLogRecord("test max level")
	assert.False(t, h.IsHandling(r.Level))
	assert.Empty(t, w.String())

	r.Level = slog.ErrorLevel
	assert.True(t, h.IsHandling(r.Level))
	assert.NoErr(t, h.Handle(r))
	assert.NotEmpty(t, w.String())

	// test handle error
	h.SetFormatter(newTestFormatter(true))
	assert.Err(t, h.Handle(r))
}
