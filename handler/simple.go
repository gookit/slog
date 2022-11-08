package handler

import (
	"io"

	"github.com/gookit/slog"
)

// SimpleHandler definition
type SimpleHandler struct {
	NopFlushClose
	slog.LevelFormattable
	Output io.Writer
}

// NewHandler create a new instance
func NewHandler(out io.Writer, level slog.Level) *SimpleHandler {
	return NewSimpleHandler(out, level)
}

// NewSimple create a new instance
func NewSimple(out io.Writer, level slog.Level) *SimpleHandler {
	return NewSimpleHandler(out, level)
}

// NewSimpleHandler create new SimpleHandler
//
// Usage:
//
//	buf := new(bytes.Buffer)
//	h := handler.NewSimpleHandler(&buf, slog.InfoLevel)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewSimpleHandler(f, slog.InfoLevel)
func NewSimpleHandler(out io.Writer, level slog.Level) *SimpleHandler {
	return &SimpleHandler{
		Output: out,
		// init log levels
		LevelFormattable: slog.NewLvFormatter(level),
	}
}

// Handle log record
func (h *SimpleHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Output.Write(bts)
	return err
}
