package slog_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gookit/rotatefile/bufwrite"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func bufferedSyncLogger(t *testing.T) (*slog.Logger, *os.File) {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), "buffered.log"))
	if err != nil {
		t.Fatal(err)
	}
	out := bufwrite.NewBufIOWriter(f)
	t.Cleanup(func() { _ = out.Close() })
	return slog.NewWithHandlers(handler.NewFlushCloseHandler(out, slog.AllLevels)), f
}

func TestLoggerSyncFlushesItsOwnHandlers(t *testing.T) {
	first, firstFile := bufferedSyncLogger(t)
	second, secondFile := bufferedSyncLogger(t)
	first.Info("first-instance-pending")
	second.Info("second-instance-pending")
	for _, f := range []*os.File{firstFile, secondFile} {
		if info, err := f.Stat(); err != nil || info.Size() != 0 {
			t.Fatalf("log did not remain buffered before Sync: info=%v, err=%v", info, err)
		}
	}

	if err := first.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(firstFile.Name())
	if err != nil || !strings.Contains(string(data), "first-instance-pending") {
		t.Errorf("Sync did not flush its receiver: data=%q, err=%v", data, err)
	}
	if info, err := secondFile.Stat(); err != nil || info.Size() != 0 {
		t.Fatalf("Sync flushed another logger: info=%v, err=%v", info, err)
	}
	if err := second.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(secondFile.Name())
	if err != nil || !strings.Contains(string(data), "second-instance-pending") {
		t.Errorf("Sync did not flush the second receiver: data=%q, err=%v", data, err)
	}
}

func TestLoggerSyncReturnsItsHandlerError(t *testing.T) {
	logger, f := bufferedSyncLogger(t)
	logger.Info("pending-before-close")
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := logger.Sync(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("Sync error = %v, want the receiver's closed-file error", err)
	}
}

type countedSyncWriter struct {
	*bufwrite.BufIOWriter
	flushes int
}

func (w *countedSyncWriter) Flush() error {
	w.flushes++
	return w.BufIOWriter.Flush()
}

func TestSugaredLoggerSyncUsesItsOwnFlushPath(t *testing.T) {
	f, err := os.Create(filepath.Join(t.TempDir(), "sugared.log"))
	if err != nil {
		t.Fatal(err)
	}
	out := &countedSyncWriter{BufIOWriter: bufwrite.NewBufIOWriter(f)}
	t.Cleanup(func() { _ = out.Close() })
	logger := slog.NewSugaredLogger(io.Discard, slog.DebugLevel)
	logger.AddHandler(handler.NewFlushCloseHandler(out, slog.AllLevels))
	logger.Info("sugared-instance-pending")
	if err := logger.Sync(); err != nil {
		t.Fatal(err)
	}
	if out.flushes != 1 {
		t.Fatalf("Sync called the attached handler %d times, want once", out.flushes)
	}
	data, err := os.ReadFile(f.Name())
	if err != nil || !strings.Contains(string(data), "sugared-instance-pending") {
		t.Fatalf("Sync did not flush the attached file: data=%q, err=%v", data, err)
	}
	if err := logger.Flush(); err != nil || out.flushes != 2 {
		t.Fatalf("Flush control: calls=%d, err=%v", out.flushes, err)
	}
}
