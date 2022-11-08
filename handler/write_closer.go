package handler

import (
	"io"

	"github.com/gookit/slog"
)

// WriteCloserHandler definition
type WriteCloserHandler struct {
	// LockWrapper
	slog.LevelFormattable
	Output io.WriteCloser
}

// NewWriteCloser create a new instance
func NewWriteCloser(out io.WriteCloser, levels []slog.Level) *WriteCloserHandler {
	return NewIOWriteCloserHandler(out, levels)
}

// NewIOWriteCloserHandler create new WriteCloserHandler
//
// Usage:
//
//	buf := new(bytes.Buffer)
//	h := handler.NewIOWriteCloserHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewIOWriteCloserHandler(f, slog.AllLevels)
func NewIOWriteCloserHandler(out io.WriteCloser, levels []slog.Level) *WriteCloserHandler {
	return &WriteCloserHandler{
		Output: out,
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// Close the handler
func (h *WriteCloserHandler) Close() error {
	return h.Output.Close()
}

// Flush the handler
func (h *WriteCloserHandler) Flush() error {
	return nil
}

// Handle log record
func (h *WriteCloserHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	// h.Lock()
	// defer h.Unlock()

	_, err = h.Output.Write(bts)
	return err
}
