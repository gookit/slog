package handler

import (
	"bufio"

	"github.com/gookit/slog"
)

// BufferedHandler definition
type BufferedHandler struct {
	lockWrapper
	LevelsWithFormatter

	buffer  *bufio.Writer
	fcWriter slog.FlushCloseWriter

	// BuffSize the buffer contents size
	// BuffSize int
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(fcWriter slog.FlushCloseWriter, bufSize int) *BufferedHandler {
	return &BufferedHandler{
		buffer:  bufio.NewWriterSize(fcWriter.Writer(), bufSize),
		fcWriter: fcWriter,
		// options
		// BuffSize: bufSize,
		// log levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
	}
}

// Flush all buffers to the `h.fcWriter.Writer()`
func (h *BufferedHandler) Flush() error {
	h.Lock()
	defer h.Unlock()

	if err := h.buffer.Flush(); err != nil {
		return err
	}

	return h.fcWriter.Flush()
}

// Close log records
func (h *BufferedHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.fcWriter.Close()
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
