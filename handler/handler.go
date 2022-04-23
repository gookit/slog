// Package handler provide useful common log handlers.
//
// eg: file, console, multi_file, rotate_file, stream, syslog, email
package handler

import (
	"bufio"
	"io"
	"os"

	"github.com/gookit/slog"
)

// defaultBufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const defaultBufferSize = 256 * 1024

var (
	// DefaultMaxSize is the maximum size of a log file in bytes.
	DefaultMaxSize uint64 = 1024 * 1024 * 1800
	// DefaultFilePerm perm and flags for create log file
	DefaultFilePerm  = 0664
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

// Builder struct for create handler
type Builder struct {
	Output   io.Writer
	Filepath string
	BuffSize int
	Levels   []slog.Level
}

// NewBuilder create
func NewBuilder() *Builder {
	return &Builder{}
}

// Build slog handler.
func (b *Builder) reset() {
	b.Output = nil
	b.Levels = b.Levels[:0]
	b.Filepath = ""
	b.BuffSize = 0
}

// Build slog handler.
func (b *Builder) Build() slog.Handler {
	defer b.reset()

	if b.Output != nil {
		return b.buildFromWriter(b.Output)
	}

	if b.Filepath != "" {
		f, err := QuickOpenFile(b.Filepath)
		if err != nil {
			panic(err)
		}

		return b.buildFromWriter(f)
	}

	panic("missing some information for build handler")
}

// Build slog handler.
func (b *Builder) buildFromWriter(w io.Writer) slog.Handler {
	if scw, ok := w.(SyncCloseWriter); ok {
		return NewSyncCloseHandler(scw, b.Levels)
	}

	if fcw, ok := w.(FlushCloseWriter); ok {
		return NewFlushCloseHandler(fcw, b.Levels)
	}

	if wc, ok := w.(io.WriteCloser); ok {
		return NewWriteCloser(wc, b.Levels)
	}

	return NewIOWriter(w, b.Levels)
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
func (h *bufFileWrapper) init() {
	if h.BuffSize > 0 {
		// TODO create buff io
	}
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
	// Level for log message. if current level <= Level will log message
	Level slog.Level
}

// create new instance
func newLvFormatter(lv slog.Level) LevelWithFormatter {
	return LevelWithFormatter{Level: lv}
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level slog.Level) bool {
	return h.Level.ShouldHandling(level)
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
