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


