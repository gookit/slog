package handler

import (
	"bufio"
	"io"

	"github.com/gookit/slog"
)

// BufferedHandler definition
type BufferedHandler struct {
	LevelsWithFormatter
	buffer *bufio.Writer
	writer io.WriteCloser
}

// NewBuffered create new BufferedHandler
func NewBuffered(cWriter io.WriteCloser, bufSize int) *BufferedHandler {
	return NewBufferedHandler(cWriter, bufSize)
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(cWriter io.WriteCloser, bufSize int) *BufferedHandler {
	return &BufferedHandler{
		writer: cWriter,
		buffer: bufio.NewWriterSize(cWriter, bufSize),
		// log levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
	}
}

// Flush all buffers
func (h *BufferedHandler) Flush() error {
	return h.buffer.Flush()
}

// Close log records
func (h *BufferedHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return h.writer.Close()
}

// Handle log record
func (h *BufferedHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.buffer.Write(bts)
	return err
}
