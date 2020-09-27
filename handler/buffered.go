package handler

import (
	"bufio"
	"sync"

	"github.com/gookit/slog"
)

const defaultFlushInterval = 1000

// BufferedHandler definition
type BufferedHandler struct {
	LevelsWithFormatter

	mu sync.Mutex

	buffer  *bufio.Writer
	handler slog.WriterHandler
	// options:
	// BuffSize for buffer
	BuffSize int
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(handler slog.WriterHandler, bufSize int) *BufferedHandler {
	return &BufferedHandler{
		buffer:  bufio.NewWriterSize(handler.Writer(), bufSize),
		handler: handler,
		// options
		BuffSize: bufSize,
	}
}

// Flush all buffers to the `h.handler.Writer()`
func (h *BufferedHandler) Flush() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := h.buffer.Flush()
	if err != nil {
		return err
	}

	return h.handler.Flush()
}

// Close log records
func (h *BufferedHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.handler.Close()
}

// Handle log record
func (h *BufferedHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.buffer == nil {
		h.buffer = bufio.NewWriterSize(h.handler.Writer(), h.BuffSize)
	}

	_, err = h.buffer.Write(bts)

	// flush logs
	if h.buffer.Buffered() >= h.BuffSize {
		return h.Flush()
	}

	return err
}
