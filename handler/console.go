package handler

import (
	"os"

	"github.com/gookit/color"
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
	h := &ConsoleHandler{}
	h.Out = os.Stdout
	h.Levels = levels

	// create new formatter
	f := slog.NewTextFormatter()
	// enable color on console
	f.EnableColor = color.IsSupportColor()

	h.SetFormatter(f)
	return h
}

// SetColorTheme to the formatter
func (h *ConsoleHandler) ConfigFormatter(fn func(formatter *slog.TextFormatter)) {
	fn(h.Formatter().(*slog.TextFormatter))
}
