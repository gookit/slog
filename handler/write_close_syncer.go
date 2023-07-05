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

// NewSyncCloserWithLF create new SyncCloseHandler, with custom slog.LevelFormattable
func NewSyncCloserWithLF(out SyncCloseWriter, lf slog.LevelFormattable) *SyncCloseHandler {
	return &SyncCloseHandler{
		Output: out,
		// init formatter and level handle
		LevelFormattable: lf,
	}
}

//
// ------------- Use max log level -------------
//

// SyncCloserWithMaxLevel create new SyncCloseHandler, with max log level
func SyncCloserWithMaxLevel(out SyncCloseWriter, maxLevel slog.Level) *SyncCloseHandler {
	return NewSyncCloserWithLF(out, slog.NewLvFormatter(maxLevel))
}

//
// ------------- Use multi log levels -------------
//

// NewSyncCloser create new SyncCloseHandler, alias of NewSyncCloseHandler()
func NewSyncCloser(out SyncCloseWriter, levels []slog.Level) *SyncCloseHandler {
	return NewSyncCloseHandler(out, levels)
}

// SyncCloserWithLevels create new SyncCloseHandler, alias of NewSyncCloseHandler()
func SyncCloserWithLevels(out SyncCloseWriter, levels []slog.Level) *SyncCloseHandler {
	return NewSyncCloseHandler(out, levels)
}

// NewSyncCloseHandler create new SyncCloseHandler with limited log levels
//
// Usage:
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewSyncCloseHandler(f, slog.AllLevels)
func NewSyncCloseHandler(out SyncCloseWriter, levels []slog.Level) *SyncCloseHandler {
	return NewSyncCloserWithLF(out, slog.NewLvsFormatter(levels))
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

	_, err = h.Output.Write(bts)
	return err
}
