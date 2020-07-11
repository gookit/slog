package handler

import (
	"io"

	"github.com/gookit/slog"
)

// StreamHandler definition
type StreamHandler struct {
	// Out io.WriteCloser
	Out io.Writer
	Levels []slog.Level
}

// Close the handler
func (h *StreamHandler) Close() error {
	// return h.Out.Close()
	return nil
}

// IsHandling Check if the current level can be handling
func (h *StreamHandler) IsHandling(level slog.Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

// Handle log record
func (h *StreamHandler) Handle(record *slog.Record) error {
	h.Out.Write()
}

// HandleBatch log records
func (h *StreamHandler) HandleBatch(records []*slog.Record) error {
	panic("implement me")
}


