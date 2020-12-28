package handler

import (
	"bufio"

	"github.com/gookit/slog"
)

// MultiFileHandler definition
type MultiFileHandler struct {
	lockWrapper
	bufio *bufio.Writer

	// FileDir for save log files
	FileDir string
	// Files can use multi file for record level logs. eg:
	//  "error.log": []slog.Level{slog.Warn, slog.Error},
	//  "info.log": []slog.Level{slog.Trace, slog.Info, slog.Notice}
	// FileLevels map[string][]slog.Level
	FileLevels map[string]slog.Levels
	// NoBuffer on write log records
	NoBuffer bool
	// BuffSize for enable buffer
	BuffSize int
	// file contents max size
	MaxSize uint64
}

// NewMultiFileHandler instance
func NewMultiFileHandler() *MultiFileHandler {
	return &MultiFileHandler{}
}

// IsHandling Check if the current level can be handling
func (h *MultiFileHandler) IsHandling(level slog.Level) bool {
	for _, ls := range h.FileLevels {
		if ls.Contains(level) {
			return true
		}
	}
	return false
}

func (h *MultiFileHandler) Close() error {
	panic("implement me")
}

func (h *MultiFileHandler) Flush() error {
	panic("implement me")
}

func (h *MultiFileHandler) Handle(record *slog.Record) error {
	panic("implement me")
}
