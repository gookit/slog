package handler

import (
	"sync"

	"github.com/gookit/slog"
)

type lockWrapper struct {
	sync.Mutex
	disable bool
}

// Lock it
func (l *lockWrapper) Lock() {
	if !l.disable {
		l.Mutex.Lock()
	}
}

// Unlock it
func (l *lockWrapper) Unlock() {
	if !l.disable {
		l.Mutex.Unlock()
	}
}

// EnableLock locker
func (l *lockWrapper) EnableLock(enable bool) {
	l.disable = !enable
}

/********************************************************************************
 * Common parts for handler
 ********************************************************************************/

// LevelsWithFormatter definition
type LevelWithFormatter struct {
	slog.Formattable
	// Level for log message. if current level >= Level will log message
	Level slog.Level
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level slog.Level) bool {
	return level >= h.Level
}

// Flush logs to disk
func (h *LevelWithFormatter) Flush() error {
	return nil
}

// Close handler
func (h *LevelWithFormatter) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return nil
}

// LevelsWithFormatter definition
type LevelsWithFormatter struct {
	slog.Formattable
	// Levels for log message
	Levels []slog.Level
}

// Flush logs to disk
func (h *LevelsWithFormatter) Flush() error {
	return nil
}

// Close handler
func (h *LevelsWithFormatter) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return nil
}

// IsHandling Check if the current level can be handling
func (h *LevelsWithFormatter) IsHandling(level slog.Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

/********************************************************************************
 * Grouped Handler
 ********************************************************************************/

// GroupedHandler definition
type GroupedHandler struct {
	handlers []slog.Handler
	// Levels for log message
	Levels []slog.Level
	// IgnoreErr on handling messages
	IgnoreErr bool
}

// NewGroupedHandler create new GroupedHandler
func NewGroupedHandler(handlers []slog.Handler) *GroupedHandler {
	return &GroupedHandler{
		handlers: handlers,
	}
}

// IsHandling Check if the current level can be handling
func (h *GroupedHandler) IsHandling(level slog.Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

// Handle log record
func (h *GroupedHandler) Handle(record *slog.Record) error {
	for _, handler := range h.handlers {
		err := handler.Handle(record)

		if h.IgnoreErr == false && err != nil {
			return err
		}
	}

	return nil
}

// Close log handlers
func (h *GroupedHandler) Close() error {
	for _, handler := range h.handlers {
		err := handler.Close()

		if h.IgnoreErr == false && err != nil {
			return err
		}
	}

	return nil
}

// Flush log records
func (h *GroupedHandler) Flush() error {
	for _, handler := range h.handlers {
		err := handler.Flush()

		if h.IgnoreErr == false && err != nil {
			return err
		}
	}

	return nil
}
