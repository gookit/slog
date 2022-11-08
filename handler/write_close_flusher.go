package handler

import (
	"github.com/gookit/slog"
)

// FlushCloseHandler definition
type FlushCloseHandler struct {
	slog.LevelFormattable
	Output FlushCloseWriter
}

// NewFlushCloser create new FlushCloseHandler
func NewFlushCloser(out FlushCloseWriter, levels []slog.Level) *FlushCloseHandler {
	return NewFlushCloseHandler(out, levels)
}

// NewFlushCloseHandler create new FlushCloseHandler
//
// Usage:
//
//	buf := new(bytes.Buffer)
//	h := handler.NewFlushCloseHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewFlushCloseHandler(f, slog.AllLevels)
func NewFlushCloseHandler(out FlushCloseWriter, levels []slog.Level) *FlushCloseHandler {
	return &FlushCloseHandler{
		Output: out,
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// Close the handler
func (h *FlushCloseHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return h.Output.Close()
}

// Flush the handler
func (h *FlushCloseHandler) Flush() error {
	return h.Output.Flush()
}

// Handle log record
func (h *FlushCloseHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	// h.Lock()
	// defer h.Unlock()

	_, err = h.Output.Write(bts)
	return err
}
