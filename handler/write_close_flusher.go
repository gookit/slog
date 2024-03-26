package handler

import (
	"github.com/gookit/slog"
)

// FlushCloseHandler definition
type FlushCloseHandler struct {
	slog.LevelFormattable
	Output FlushCloseWriter
}

// NewFlushCloserWithLF create new FlushCloseHandler, with custom slog.LevelFormattable
func NewFlushCloserWithLF(out FlushCloseWriter, lf slog.LevelFormattable) *FlushCloseHandler {
	return &FlushCloseHandler{
		Output: out,
		// init formatter and level handle
		LevelFormattable: lf,
	}
}

//
// ------------- Use max log level -------------
//

// FlushCloserWithMaxLevel create new FlushCloseHandler, with max log level
func FlushCloserWithMaxLevel(out FlushCloseWriter, maxLevel slog.Level) *FlushCloseHandler {
	return NewFlushCloserWithLF(out, slog.NewLvFormatter(maxLevel))
}

//
// ------------- Use multi log levels -------------
//

// NewFlushCloser create new FlushCloseHandler, alias of NewFlushCloseHandler()
func NewFlushCloser(out FlushCloseWriter, levels []slog.Level) *FlushCloseHandler {
	return NewFlushCloseHandler(out, levels)
}

// FlushCloserWithLevels create new FlushCloseHandler, alias of NewFlushCloseHandler()
func FlushCloserWithLevels(out FlushCloseWriter, levels []slog.Level) *FlushCloseHandler {
	return NewFlushCloseHandler(out, levels)
}

// NewFlushCloseHandler create new FlushCloseHandler
//
// Usage:
//
//	buf := new(byteutil.Buffer)
//	h := handler.NewFlushCloseHandler(&buf, slog.AllLevels)
//
//	f, err := os.OpenFile("my.log", ...)
//	h := handler.NewFlushCloseHandler(f, slog.AllLevels)
func NewFlushCloseHandler(out FlushCloseWriter, levels []slog.Level) *FlushCloseHandler {
	return NewFlushCloserWithLF(out, slog.NewLvsFormatter(levels))
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

	_, err = h.Output.Write(bts)
	return err
}
