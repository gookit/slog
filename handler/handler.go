// Package handler provide some common log handlers.
// eg: file, console, multi_file, rotate_file, stream, syslog, email
package handler

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/gookit/slog"
)

// var (
// 	// program pid
// 	pid = os.Getpid()
// 	// program name
// 	pName = filepath.Base(os.Args[0])
// 	hName = "unknownHost" // TODO
// 	// uName = "unknownUser"
// )

// defaultBufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const defaultBufferSize = 256 * 1024

var (
	// DefaultMaxSize is the maximum size of a log file in bytes.
	DefaultMaxSize uint64 = 1024 * 1024 * 1800
	// perm and flags for create log file
	DefaultFilePerm  = 0664
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type lockWrapper struct {
	sync.Mutex
	disable bool
}

// Lock it
func (lw *lockWrapper) Lock() {
	if false == lw.disable {
		lw.Mutex.Lock()
	}
}

// Unlock it
func (lw *lockWrapper) Unlock() {
	if !lw.disable {
		lw.Mutex.Unlock()
	}
}

// UseLock locker
func (lw *lockWrapper) UseLock(enable bool) {
	lw.disable = false == enable
}

// LockEnabled status
func (lw *lockWrapper) LockEnabled() bool {
	return lw.disable == false
}

// NopFlushClose no operation. provide empty Flush(), Close() methods
type NopFlushClose struct{}

// Flush logs to disk
func (h *NopFlushClose) Flush() error {
	return nil
}

// Close handler
func (h *NopFlushClose) Close() error {
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

	file, err := QuickOpenFile(h.fpath)
	if err != nil {
		return err
	}

	h.file = file
	return err
}

// Write contents to *os.File
func (h *fileWrapper) Write(bts []byte) (n int, err error) {
	return h.file.Write(bts)
}

// Writer return *os.File
func (h *fileWrapper) Writer() io.Writer {
	return h.file
}

// Close handler, will be flush logs to file, then close file
func (h *fileWrapper) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.file.Close()
}

// Flush logs to disk file
func (h *fileWrapper) Flush() error {
	return h.file.Sync()
}

type bufFileWrapper struct {
	fileWrapper
	bufio *bufio.Writer

	written uint64
	// NoBuffer on write log records
	NoBuffer bool
	// BuffSize for enable buffer
	BuffSize int
}

// CloseBuffer for write logs
func (h *bufFileWrapper) CloseBuffer() {
	h.NoBuffer = true
}

// Write logs
func (h *bufFileWrapper) Write(bts []byte) (n int, err error) {
	// direct write logs to file
	if h.NoBuffer {
		n, err = h.file.Write(bts)
	} else {
		// enable buffer
		if h.bufio == nil && h.BuffSize > 0 {
			h.bufio = bufio.NewWriterSize(h.file, h.BuffSize)
		}

		n, err = h.bufio.Write(bts)
	}
	return
}

// Flush logs to disk file
func (h *bufFileWrapper) Flush() error {
	// flush buffers to h.file
	if h.bufio != nil {
		err := h.bufio.Flush()
		if err != nil {
			return err
		}
	}

	return h.file.Sync()
}

// Close handler, will be flush logs to file, then close file
func (h *bufFileWrapper) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.file.Close()
}

/********************************************************************************
 * Common parts for handler
 ********************************************************************************/

// LevelWithFormatter struct definition
//
// - support set log formatter
// - only support set one log level
type LevelWithFormatter struct {
	slog.Formattable
	// Level for log message. if current level >= Level will log message
	Level slog.Level
}

// create new instance
func newLvFormatter(lv slog.Level) LevelWithFormatter {
	return LevelWithFormatter{Level: lv}
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level slog.Level) bool {
	return level >= h.Level
}

// LevelsWithFormatter struct definition
//
// - support set log formatter
// - support setting multi log levels
type LevelsWithFormatter struct {
	slog.Formattable
	// Levels for log message
	Levels []slog.Level
}

// create new instance
func newLvsFormatter(lvs []slog.Level) LevelsWithFormatter {
	return LevelsWithFormatter{Levels: lvs}
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
func (h *GroupedHandler) Handle(record *slog.Record) (err error) {
	for _, handler := range h.handlers {
		err = handler.Handle(record)

		if h.IgnoreErr == false && err != nil {
			return err
		}
	}
	return
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
