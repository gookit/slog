package handler_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/x/fakeobj"
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

	str := fsutil.ReadString(logfile)
	assert.Contains(t, str, "test sync closer handler")
	assert.NoErr(t, h.Close())

	t.Run("err on sync", func(t *testing.T) {
		w := &syncCloseWriter{}
		w.errOnSync = true
		h = handler.SyncCloserWithLevels(w, slog.NormalLevels)

		assert.Err(t, h.Flush())
		assert.Err(t, h.Close())
	})

	// test handle error
	h.SetFormatter(newTestFormatter(true))
	assert.Err(t, h.Handle(r))
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

	t.Run("use max level", func(t *testing.T) {
		h = handler.WriteCloserWithMaxLevel(w, slog.WarnLevel)
		r = newLogRecord("test max level")
		assert.False(t, h.IsHandling(r.Level))

		r.Level = slog.ErrorLevel
		assert.True(t, h.IsHandling(r.Level))
	})

	// test handle error
	t.Run("handle error", func(t *testing.T) {
		h = handler.WriteCloserWithLevels(w, slog.NormalLevels)
		h.SetFormatter(newTestFormatter(true))
		assert.Err(t, h.Handle(r))
	})
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

	t.Run("ErrOnFlush", func(t *testing.T) {
		w.ErrOnFlush = true
		assert.Err(t, h.Flush())
		assert.Err(t, h.Close())
	})

	t.Run("With max level", func(t *testing.T) {
		h = handler.FlushCloserWithMaxLevel(w, slog.WarnLevel)
		r = newLogRecord("test max level")
		assert.False(t, h.IsHandling(r.Level))
		assert.Empty(t, w.String())

		r.Level = slog.ErrorLevel
		assert.True(t, h.IsHandling(r.Level))
		assert.NoErr(t, h.Handle(r))
		assert.NotEmpty(t, w.String())
	})

	// test handle error
	h = handler.FlushCloserWithMaxLevel(w, slog.WarnLevel)
	h.SetFormatter(newTestFormatter(true))
	assert.Err(t, h.Handle(r))
}

func TestNewSimpleHandler(t *testing.T) {
	buf := fakeobj.NewWriter()

	h := handler.NewSimple(buf, slog.InfoLevel)
	r := newLogRecord("test simple handler")
	assert.NoErr(t, h.Handle(r))

	s := buf.String()
	buf.Reset()
	fmt.Print(s)
	assert.Contains(t, s, "test simple handler")

	assert.NoErr(t, h.Flush())
	assert.NoErr(t, h.Close())

	h = handler.NewHandler(buf, slog.InfoLevel)
	r = newLogRecord("test simple handler2")
	assert.NoErr(t, h.Handle(r))

	s = buf.ResetGet()
	fmt.Print(s)
	assert.Contains(t, s, "test simple handler2")

	assert.NoErr(t, h.Flush())
	assert.NoErr(t, h.Close())

	h = handler.SimpleWithLevels(buf, slog.NormalLevels)
	r = newLogRecord("test simple handler with levels")
	assert.NoErr(t, h.Handle(r))

	s = buf.ResetGet()
	fmt.Print(s)
	assert.Contains(t, s, "test simple handler with levels")

	// handle error
	h.SetFormatter(newTestFormatter(true))
	assert.Err(t, h.Handle(r))
}
