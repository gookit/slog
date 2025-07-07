package slog

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gookit/goutil/strutil"
)

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
	// Level max for logging messages. if current level <= Level will log messages
	Level Level
}

// NewLvFormatter create new LevelWithFormatter instance
func NewLvFormatter(maxLv Level) *LevelWithFormatter {
	return &LevelWithFormatter{Level: maxLv}
}

// SetMaxLevel set max level for logging messages
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
	// Levels for logging messages
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

// LevelMode define level mode for logging
type LevelMode uint8

// MarshalJSON implement the JSON Marshal interface [encoding/json.Marshaler]
func (m LevelMode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.String() + `"`), nil
}

// UnmarshalJSON implement the JSON Unmarshal interface [encoding/json.Unmarshaler]
func (m *LevelMode) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	*m, err = StringToLevelMode(s)
	return err
}

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

// SafeToLevelMode parse string value to LevelMode, fail return LevelModeList
func SafeToLevelMode(s string) LevelMode {
	lm, err := StringToLevelMode(s)
	if err != nil {
		return LevelModeList
	}
	return lm
}

// StringToLevelMode parse string value to LevelMode
func StringToLevelMode(s string) (LevelMode, error) {
	switch s {
	case "", "list", "list_level", "level_list":
		return LevelModeList, nil
	case "max", "max_level", "level_max":
		return LevelModeMax, nil
	default:
		// is int value, try to parse as int
		if strutil.IsInt(s) {
			iVal := strutil.SafeInt(s)
			if iVal >= 0 && iVal <= int(LevelModeMax) {
				return LevelMode(iVal), nil
			}
		}
		return 0, fmt.Errorf("slog: invalid level mode: %s", s)
	}
}

// LevelHandling struct definition
type LevelHandling struct {
	// level check mode. default is LevelModeList
	lvMode LevelMode
	// max level for a log message. if the current level <= Level will log a message
	maxLevel Level
	// levels limit for log message
	levels []Level
}

// SetMaxLevel set max level for a log message
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
