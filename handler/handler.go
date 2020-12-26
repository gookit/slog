// Package handler provide some common log handlers.
// eg: file, console, multi_file, rotate_file, stream, syslog, email
package handler

import (
	"io"
	"os"
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

// UseLock locker
func (l *lockWrapper) UseLock(enable bool) {
	l.disable = false == enable
}

// LockEnabled status
func (l *lockWrapper) LockEnabled() bool {
	return l.disable == false
}

type emptyHandler struct {}

// Flush logs to disk
func (h *emptyHandler) Flush() error {
	return nil
}

// Close handler
func (h *emptyHandler) Close() error {
	return nil
}

type fileWrapper struct {
	fpath string
	file  *os.File
}

// func newFileHandler(fpath string) *fileWrapper {
// 	return &fileWrapper{fpath: fpath}
// }

// ReopenFile the log file
func (h *fileWrapper) ReopenFile() error {
	if h.file != nil {
		h.file.Close()
	}

	file, err := openFile(h.fpath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return err
	}

	h.file = file
	return err
}

// Writer return *os.File
func (h *fileWrapper) Writer() io.Writer {
	return h.file
}

// Close handler, will be flush logs to file, then close file
func (h *fileWrapper) Close() error {
	if err := h.file.Sync(); err != nil {
		return err
	}

	return h.file.Close()
}

// Flush logs to disk file
func (h *fileWrapper) Flush() error {
	return h.file.Sync()
}

/********************************************************************************
 * Common parts for handler
 ********************************************************************************/

// LevelsWithFormatter struct definition
//
// - support set log formatter
// - support setting multi log levels
type LevelWithFormatter struct {
	slog.Formattable
	// Level for log message. if current level >= Level will log message
	Level slog.Level
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level slog.Level) bool {
	return level >= h.Level
}

// LevelsWithFormatter struct definition
//
// - support set log formatter
// - only support set one log level
type LevelsWithFormatter struct {
	slog.Formattable
	// Levels for log message
	Levels []slog.Level
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
