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

// Enable locker
func (l *lockWrapper) Enable(enable bool) {
	l.disable = !enable
}

/********************************************************************************
 * Base handler
 ********************************************************************************/

// BaseHandler definition
type BaseHandler struct {
	slog.Formattable
	// Levels for log
	Levels []slog.Level
}

func (h *BaseHandler) Flush() error {
	return nil
}

// IsHandling Check if the current level can be handling
func (h *BaseHandler) IsHandling(level slog.Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}

// HandleBatch log records
func (h *BaseHandler) HandleBatch(records []*slog.Record) error {
	panic("implement me")
}

/********************************************************************************
 * Grouped Handler
 ********************************************************************************/

// GroupedHandler definition
type GroupedHandler struct {
	handlers []slog.Handler
	// Levels for log
	Levels    []slog.Level
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
