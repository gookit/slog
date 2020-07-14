package handler

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/formatter"
)

/********************************************************************************
 * Base handler
 ********************************************************************************/

// BaseHandler definition
type BaseHandler struct {
	formatter.Formattable
	// Levels for log
	Levels []slog.Level
}

func (h *BaseHandler) Flush() error  {
	return nil
}

// IsHandling Check if the current level can be handling
func (h *BaseHandler) IsHandling(level slog.Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

// HandleBatch log records
func (h *BaseHandler) HandleBatch(records []*slog.Record) error {
	panic("implement me")
}
