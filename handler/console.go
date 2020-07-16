package handler

import (
	"os"

	"github.com/gookit/slog"
)

/********************************************************************************
 * console log handler
 ********************************************************************************/

// ConsoleHandler definition
type ConsoleHandler struct {
	StreamHandler
}

// NewConsoleHandler create new ConsoleHandler
func NewConsoleHandler(levels []slog.Level) *ConsoleHandler {
	h := &ConsoleHandler{
		StreamHandler: *NewStreamHandler(os.Stdout, levels),
	}

	// create new formatter
	f := slog.NewTextFormatter()
	// enable color
	f.EnableColor = true

	h.SetFormatter(f)
	return h
}
