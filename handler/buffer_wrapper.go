package handler

import (
	"io"

	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
)

// FormatterWriterHandler interface
type FormatterWriterHandler interface {
	slog.Handler
	// Formatter record formatter
	Formatter() slog.Formatter
	// Writer the output writer
	Writer() io.Writer
}

// bufferWrapper struct
type bufferWrapper struct {
	buffer  FlushWriter
	handler FormatterWriterHandler
}

// BufferWrapper new instance.
func BufferWrapper(handler FormatterWriterHandler, buffSize int) slog.Handler {
	return &bufferWrapper{
		handler: handler,
		buffer:  bufwrite.NewBufIOWriterSize(handler.Writer(), buffSize),
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
