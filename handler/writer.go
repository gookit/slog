package handler

import (
	"io"

	"github.com/gookit/slog"
)

// IOWriterHandler definition
type IOWriterHandler struct {
	slog.LevelFormattable
	Output io.Writer
}

// NewIOWriter create a new instance
func NewIOWriter(out io.Writer, levels []slog.Level) *IOWriterHandler {
	return NewIOWriterHandler(out, levels)
}

// NewIOWriterHandler create new IOWriterHandler
//
// Usage:
//
//	buf := new(bytes.Buffer)
//	h := handler.NewIOWriterHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewIOWriterHandler(f, slog.AllLevels)
func NewIOWriterHandler(out io.Writer, levels []slog.Level) *IOWriterHandler {
	return &IOWriterHandler{
		Output: out,
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// Close the handler
func (h *IOWriterHandler) Close() error {
	return nil
}

// Flush the handler
func (h *IOWriterHandler) Flush() error {
	return nil
}

// Handle log record
func (h *IOWriterHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Output.Write(bts)
	return err
}
