package handler

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/formatter"
)

var defaultFormatter = formatter.LineFormatter{}

// BaseHandler definition
type BaseHandler struct {

}

func (h *BaseHandler) Flush() error  {
	return nil
}

// HandleBatch log records
func (h *BaseHandler) HandleBatch(records []*slog.Record) error {
	panic("implement me")
}

// BufferedHandler definition
type BufferedHandler struct {

}
