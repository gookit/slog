package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewIOWriter(t *testing.T) {
	w := new(bytes.Buffer)
	h := handler.NewIOWriter(w, slog.NormalLevels)

	assert.True(t, h.IsHandling(slog.NoticeLevel))

	r := newLogRecord("test io.writer handler")
	assert.NoError(t, h.Handle(r))
	assert.NoError(t, h.Flush())

	str := w.String()
	assert.Contains(t, str, "test io.writer handler")

	assert.NoError(t, h.Close())
}

func TestNewSyncCloser(t *testing.T) {
	logfile := "./testdata/sync-closer.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	f, err := handler.QuickOpenFile(logfile)
	assert.NoError(t, err)

	h := handler.NewSyncCloser(f, slog.DangerLevels)

	assert.True(t, h.IsHandling(slog.WarnLevel))
	assert.False(t, h.IsHandling(slog.InfoLevel))

	r := newLogRecord("test sync closer handler")
	r.Level = slog.ErrorLevel

	err = h.Handle(r)
	assert.NoError(t, err)
	assert.NoError(t, h.Flush())

	bts := fsutil.MustReadFile(logfile)
	str := string(bts)

	assert.Contains(t, str, "test sync closer handler")

	err = h.Close()
	assert.NoError(t, err)
}

type closedBuffer struct {
	bytes.Buffer
}

func (w *closedBuffer) Close() error {
	return nil
}

func TestNewWriteCloser(t *testing.T) {
	w := &closedBuffer{Buffer: bytes.Buffer{}}
	h := handler.NewWriteCloser(w, slog.NormalLevels)

	assert.True(t, h.IsHandling(slog.NoticeLevel))

	r := newLogRecord("test writeCloser handler")
	assert.NoError(t, h.Handle(r))
	assert.NoError(t, h.Flush())

	str := w.String()
	assert.Contains(t, str, "test writeCloser handler")

	assert.NoError(t, h.Close())
}
