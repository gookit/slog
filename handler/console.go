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
type ConsoleHandler = IOWriterHandler

// NewConsoleWithLF create new ConsoleHandler and with custom slog.LevelFormattable
func NewConsoleWithLF(lf slog.LevelFormattable) *ConsoleHandler {
	h := NewIOWriterWithLF(os.Stdout, lf)

	// default use text formatter
	f := slog.NewTextFormatter()
	// default enable color on console
	f.WithEnableColor(color.SupportColor())

	h.SetFormatter(f)
	return h
}

//
// ------------- Use max log level -------------
//

// ConsoleWithMaxLevel create new ConsoleHandler and with max log level
func ConsoleWithMaxLevel(level slog.Level) *ConsoleHandler {
	return NewConsoleWithLF(slog.NewLvFormatter(level))
}

//
// ------------- Use multi log levels -------------
//

// NewConsole create new ConsoleHandler, alias of NewConsoleHandler
func NewConsole(levels []slog.Level) *ConsoleHandler {
	return NewConsoleHandler(levels)
}

// ConsoleWithLevels create new ConsoleHandler and with limited log levels
func ConsoleWithLevels(levels []slog.Level) *ConsoleHandler {
	return NewConsoleHandler(levels)
}

// NewConsoleHandler create new ConsoleHandler with limited log levels
func NewConsoleHandler(levels []slog.Level) *ConsoleHandler {
	return NewConsoleWithLF(slog.NewLvsFormatter(levels))
}
