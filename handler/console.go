package handler

import (
	"github.com/gookit/color"
	"github.com/gookit/slog"
)

/*************************************************************
 * console log
 *************************************************************/

var ColorTheme = map[slog.Level]color.Color{
	slog.ErrorLevel: color.FgRed,
	slog.WarnLevel:  color.FgYellow,
	slog.InfoLevel:  color.FgGreen,
	slog.DebugLevel: color.FgCyan,
	slog.TraceLevel: color.FgMagenta,
}

// ConsoleHandler definition
type ConsoleHandler struct {

}

func (c *ConsoleHandler) Close() error {
	panic("implement me")
}

func (c *ConsoleHandler) IsHandling(level slog.Level) bool {
	panic("implement me")
}

func (c *ConsoleHandler) Handle(record *slog.Record) bool {
	panic("implement me")
}

func (c *ConsoleHandler) HandleBatch(records []*slog.Record) bool {
	panic("implement me")
}
