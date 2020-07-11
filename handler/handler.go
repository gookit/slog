package handler

import "github.com/gookit/slog/formatter"

var defaultFormatter = formatter.LineFormatter{}

// BaseHandler definition
type BaseHandler struct {

}

func (h *BaseHandler) HandleBatch()  {

}

// BufferedHandler definition
type BufferedHandler struct {

}
