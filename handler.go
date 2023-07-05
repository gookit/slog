package slog

import "io"

//
// Handler interface
//

// Handler interface definition
type Handler interface {
	// Closer Close handler.
	// You should first call Flush() on close logic.
	// Refer the FileHandler.Close() handle
	io.Closer
	// Flush and sync logs to disk file.
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	//
	// All records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
}

// LevelFormattable support limit log levels and provide formatter
type LevelFormattable interface {
	Formattable
	IsHandling(level Level) bool
}

// FormattableHandler interface
type FormattableHandler interface {
	Handler
	Formattable
}

/********************************************************************************
 * Common parts for handler
 ********************************************************************************/

// LevelWithFormatter struct definition
//
// - support set log formatter
// - only support set max log level
type LevelWithFormatter struct {
	FormattableTrait
	// Level max for log message. if current level <= Level will log message
	Level Level
}

// NewLvFormatter create new LevelWithFormatter instance
func NewLvFormatter(maxLv Level) *LevelWithFormatter {
	return &LevelWithFormatter{Level: maxLv}
}

// SetMaxLevel set max level for log message
func (h *LevelWithFormatter) SetMaxLevel(maxLv Level) {
	h.Level = maxLv
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level Level) bool {
	return h.Level.ShouldHandling(level)
}

// LevelsWithFormatter struct definition
//
// - support set log formatter
// - support setting multi log levels
type LevelsWithFormatter struct {
	FormattableTrait
	// Levels for log message
	Levels []Level
}

// NewLvsFormatter create new instance
func NewLvsFormatter(levels []Level) *LevelsWithFormatter {
	return &LevelsWithFormatter{Levels: levels}
}

// SetLimitLevels set limit levels for log message
func (h *LevelsWithFormatter) SetLimitLevels(levels []Level) {
	h.Levels = levels
}

// IsHandling Check if the current level can be handling
func (h *LevelsWithFormatter) IsHandling(level Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

// LevelMode define level mode
type LevelMode uint8

// String return string value
func (m LevelMode) String() string {
	switch m {
	case LevelModeList:
		return "list"
	case LevelModeMax:
		return "max"
	default:
		return "unknown"
	}
}

const (
	// LevelModeList use level list for limit record write
	LevelModeList LevelMode = iota
	// LevelModeMax use max level limit log record write
	LevelModeMax
)

// LevelHandling struct definition
type LevelHandling struct {
	// level check mode. default is LevelModeList
	lvMode LevelMode
	// max level for log message. if current level <= Level will log message
	maxLevel Level
	// levels limit for log message
	levels []Level
}

// SetMaxLevel set max level for log message
func (h *LevelHandling) SetMaxLevel(maxLv Level) {
	h.lvMode = LevelModeMax
	h.maxLevel = maxLv
}

// SetLimitLevels set limit levels for log message
func (h *LevelHandling) SetLimitLevels(levels []Level) {
	h.lvMode = LevelModeList
	h.levels = levels
}

// IsHandling Check if the current level can be handling
func (h *LevelHandling) IsHandling(level Level) bool {
	if h.lvMode == LevelModeMax {
		return h.maxLevel.ShouldHandling(level)
	}

	for _, l := range h.levels {
		if l == level {
			return true
		}
	}
	return false
}

// LevelFormatting wrap level handling and log formatter
type LevelFormatting struct {
	LevelHandling
	FormatterWrapper
}

// NewMaxLevelFormatting create new instance with max level
func NewMaxLevelFormatting(maxLevel Level) *LevelFormatting {
	lf := &LevelFormatting{}
	lf.SetMaxLevel(maxLevel)
	return lf
}

// NewLevelsFormatting create new instance with levels
func NewLevelsFormatting(levels []Level) *LevelFormatting {
	lf := &LevelFormatting{}
	lf.SetLimitLevels(levels)
	return lf
}
