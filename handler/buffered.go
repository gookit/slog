package handler

import (
	"bufio"

	"github.com/gookit/slog"
)

const defaultFlushInterval = 1000

// BufferedHandler definition
type BufferedHandler struct {
	BaseHandler
	bufio.Writer
	handler slog.Handler
	number int
	// options
	FlushInterval int
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(handler slog.Handler) *BufferedHandler {
	return &BufferedHandler{
		handler: handler,
		FlushInterval: defaultFlushInterval,
	}
}

func (h *BufferedHandler) Flush() error {
	return h.Writer.Flush()
}

// Handle log record
func (h *BufferedHandler) Handle(record *slog.Record)  error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	h.number++
	_, err = h.Writer.Write(bts)

	if h.number >= h.FlushInterval {
		return h.Flush()
	}

	return err
}

