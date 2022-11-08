package handler

import (
	"io"

	"github.com/gookit/slog"
)

// SyncCloseHandler definition
type SyncCloseHandler struct {
	slog.LevelFormattable
	Output SyncCloseWriter
}

// NewSyncCloser create new SyncCloseHandler
func NewSyncCloser(out SyncCloseWriter, levels []slog.Level) *SyncCloseHandler {
	return NewSyncCloseHandler(out, levels)
}

// NewSyncCloseHandler create new SyncCloseHandler
//
// Usage:
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewSyncCloseHandler(f, slog.AllLevels)
func NewSyncCloseHandler(out SyncCloseWriter, levels []slog.Level) *SyncCloseHandler {
	return &SyncCloseHandler{
		Output: out,
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// Close the handler
func (h *SyncCloseHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return h.Output.Close()
}

// Flush the handler
func (h *SyncCloseHandler) Flush() error {
	return h.Output.Sync()
}

// Writer of the handler
func (h *SyncCloseHandler) Writer() io.Writer {
	return h.Output
}

// Handle log record
func (h *SyncCloseHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	// h.Lock()
	// defer h.Unlock()

	_, err = h.Output.Write(bts)
	return err
}
