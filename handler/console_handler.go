package handler

import (
	"os"

	"github.com/gookit/color"
	"github.com/gookit/slog"
)

/*************************************************************
 * console log
 *************************************************************/

// ColorTheme for console log
var ColorTheme = map[slog.Level]color.Color{
	slog.ErrorLevel: color.FgRed,
	slog.WarnLevel:  color.FgYellow,
	slog.InfoLevel:  color.FgGreen,
	slog.DebugLevel: color.FgCyan,
	slog.TraceLevel: color.FgMagenta,
}

// ConsoleHandler definition
type ConsoleHandler struct {
	StreamHandler
}

// NewConsoleHandler create new ConsoleHandler
func NewConsoleHandler(levels []slog.Level) *ConsoleHandler {
	return &ConsoleHandler{
		StreamHandler{
			Out: os.Stdout,
			Levels: levels,
		},
	}
}

