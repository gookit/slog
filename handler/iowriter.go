package handler

import (
	"io"

	"github.com/gookit/slog"
)

// IOWriterHandler definition
type IOWriterHandler struct {
	lockWrapper
	LevelsWithFormatter

	// Output io.WriteCloser
	Output io.Writer
}

// NewIOWriter create an new instance
func NewIOWriter(out io.Writer, levels []slog.Level) *IOWriterHandler {
	return NewIOWriterHandler(out, levels)
}

// NewIOWriterHandler create new IOWriterHandler
// Usage:
// 	buf := new(bytes.Buffer)
// 	h := handler.NewIOWriterHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
// 	h := handler.NewIOWriterHandler(f, slog.AllLevels)
func NewIOWriterHandler(out io.Writer, levels []slog.Level) *IOWriterHandler {
	return &IOWriterHandler{
		Output: out,
		// init log levels
		LevelsWithFormatter: newLvsFormatter(levels),
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

	h.Lock()
	defer h.Unlock()

	_, err = h.Output.Write(bts)
	return err
}
