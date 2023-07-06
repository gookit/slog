package handler

import (
	"io"

	"github.com/gookit/slog"
)

// WriteCloserHandler definition
type WriteCloserHandler struct {
	slog.LevelFormattable
	Output io.WriteCloser
}

// NewWriteCloserWithLF create new WriteCloserHandler and with custom slog.LevelFormattable
func NewWriteCloserWithLF(out io.WriteCloser, lf slog.LevelFormattable) *WriteCloserHandler {
	return &WriteCloserHandler{
		Output: out,
		// init formatter and level handle
		LevelFormattable: lf,
	}
}

// WriteCloserWithMaxLevel create new WriteCloserHandler and with max log level
func WriteCloserWithMaxLevel(out io.WriteCloser, maxLevel slog.Level) *WriteCloserHandler {
	return NewWriteCloserWithLF(out, slog.NewLvFormatter(maxLevel))
}

//
// ------------- Use multi log levels -------------
//

// WriteCloserWithLevels create a new instance and with limited log levels
func WriteCloserWithLevels(out io.WriteCloser, levels []slog.Level) *WriteCloserHandler {
	// h := &WriteCloserHandler{Output: out}
	// h.LimitLevels(levels)
	return NewWriteCloserHandler(out, levels)
}

// NewWriteCloser create a new instance
func NewWriteCloser(out io.WriteCloser, levels []slog.Level) *WriteCloserHandler {
	return NewWriteCloserHandler(out, levels)
}

// NewWriteCloserHandler create new WriteCloserHandler
//
// Usage:
//
//	buf := new(bytes.Buffer)
//	h := handler.NewIOWriteCloserHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewIOWriteCloserHandler(f, slog.AllLevels)
func NewWriteCloserHandler(out io.WriteCloser, levels []slog.Level) *WriteCloserHandler {
	return NewWriteCloserWithLF(out, slog.NewLvsFormatter(levels))
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

	_, err = h.Output.Write(bts)
	return err
}
