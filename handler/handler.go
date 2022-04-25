// Package handler provide useful common log handlers.
//
// eg: file, console, multi_file, rotate_file, stream, syslog, email
package handler

import (
	"io"
	"os"

	"github.com/gookit/slog"
)

// DefaultBufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
var DefaultBufferSize = 8 * 1024

var (
	// DefaultFilePerm perm and flags for create log file
	DefaultFilePerm os.FileMode = 0664
	// DefaultFileFlags for create/open file
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
