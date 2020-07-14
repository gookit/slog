package handler

import (
	"io"

	"github.com/gookit/slog"
)

// StreamHandler definition
type StreamHandler struct {
	BaseHandler

	// Out io.WriteCloser
	Out io.Writer

	FilePerm int
	UseLock bool
}

// NewStreamHandler create new StreamHandler
func NewStreamHandler(out io.Writer, levels []slog.Level) *StreamHandler {
	return &StreamHandler{
		Out: out,
		BaseHandler: BaseHandler{
			Levels: levels,
		},
	}
}

// Close the handler
func (h *StreamHandler) Close() error {
	// return h.Out.Close()
	return nil
}

// Handle log record
func (h *StreamHandler) Handle(record *slog.Record)  error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Out.Write(bts)
	return err
}



