package handler

import (
	"bufio"
	"io"

	"github.com/gookit/slog"
)

// BufferedHandler definition
type BufferedHandler struct {
	lockWrapper
	LevelsWithFormatter

	buffer  *bufio.Writer
	cWriter io.WriteCloser
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(cWriter io.WriteCloser, bufSize int) *BufferedHandler {
	return &BufferedHandler{
		cWriter: cWriter,
		buffer:  bufio.NewWriterSize(cWriter, bufSize),
		// log levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
	}
}

// Flush all buffers
func (h *BufferedHandler) Flush() error {
	h.Lock()
	defer h.Unlock()

	return h.buffer.Flush()
}

// Close log records
func (h *BufferedHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.cWriter.Close()
}

// Handle log record
func (h *BufferedHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()

	// if h.buffer == nil {
	// 	h.buffer = bufio.NewWriterSize(h.fcWriter.Writer(), h.BuffSize)
	// }

	_, err = h.buffer.Write(bts)
	return err
}
