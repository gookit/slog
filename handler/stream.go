package handler

import (
	"io"

	"github.com/gookit/slog"
)

// StreamHandler definition
type StreamHandler struct {
	BaseHandler

	// Output io.WriteCloser
	Output  io.Writer
	UseLock bool
}

// NewStreamHandler create new StreamHandler
func NewStreamHandler(out io.Writer, levels []slog.Level) *StreamHandler {
	return &StreamHandler{
		Output: out,
		BaseHandler: BaseHandler{
			Levels: levels,
		},
	}
}

// Close the handler
func (h *StreamHandler) Close() error {
	return h.Flush()
}

// Handle log record
func (h *StreamHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Output.Write(bts)
	return err
}
