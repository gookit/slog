package handler

import (
	"bufio"

	"github.com/gookit/slog"
)

// bufferWrapper struct
type bufferWrapper struct {
	buffer  *bufio.Writer
	handler slog.FormatterWriterHandler
}

// BufferWrapper new instance
func BufferWrapper(handler slog.FormatterWriterHandler, buffSize int) *bufferWrapper {
	return &bufferWrapper{
		handler: handler,
		buffer:  bufio.NewWriterSize(handler.Writer(), buffSize),
	}
}

// IsHandling Check if the current level can be handling
func (w *bufferWrapper) IsHandling(level slog.Level) bool {
	return w.handler.IsHandling(level)
}

// Flush all buffers to the `h.fcWriter.Writer()`
func (w *bufferWrapper) Flush() error {
	if err := w.buffer.Flush(); err != nil {
		return err
	}

	return w.handler.Flush()
}

// Close log records
func (w *bufferWrapper) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}

	return w.handler.Close()
}

// Handle log record
func (w *bufferWrapper) Handle(record *slog.Record) error {
	bts, err := w.handler.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = w.buffer.Write(bts)
	return err
}
