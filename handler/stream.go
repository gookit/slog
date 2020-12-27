package handler

import (
	"io"

	"github.com/gookit/slog"
)

// StreamHandler definition
type StreamHandler struct {
	lockWrapper
	LevelsWithFormatter

	// Output io.WriteCloser
	Output io.Writer
}

// NewStreamHandler create new StreamHandler
// Usage:
// 	buf := new(bytes.Buffer)
// 	h := handler.NewStreamHandler(&buf, slog.AllLevels)
//	f, err := os.OpenFile("my.log", ...)
// 	h := handler.NewStreamHandler(f, slog.AllLevels)
func NewStreamHandler(out io.Writer, levels []slog.Level) *StreamHandler {
	return &StreamHandler{
		Output: out,
		// init log levels
		LevelsWithFormatter: LevelsWithFormatter{
			Levels: levels,
		},
	}
}

// Close the handler
func (h *StreamHandler) Close() error {
	return nil
}

// Flush the handler
func (h *StreamHandler) Flush() error {
	return nil
}

// Handle log record
func (h *StreamHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()

	_, err = h.Output.Write(bts)
	return err
}
