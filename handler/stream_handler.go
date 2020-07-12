package handler

import (
	"io"

	"github.com/gookit/slog"
	"github.com/gookit/slog/formatter"
)

// StreamHandler definition
type StreamHandler struct {
	formatter.Formattable
	// Out io.WriteCloser
	Out io.Writer

	Levels []slog.Level

	FilePerm int
	UseLock bool
}

func NewStreamHandler(out io.Writer, levels []slog.Level) *StreamHandler {
	return &StreamHandler{
		Out: out,
		Levels: levels,
	}
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
func (h *StreamHandler) Handle(record *slog.Record) (err error) {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Out.Write(bts)
	return
}

// HandleBatch log records
func (h *StreamHandler) HandleBatch(records []*slog.Record) error {
	panic("implement me")
}


