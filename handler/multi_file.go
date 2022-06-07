package handler

import (
	"github.com/gookit/slog"
)

// MultiFileHandler definition TODO
type MultiFileHandler struct {
	LockWrapper
	// writers map[string]io.Writer
	// FileDir for save log files
	FileDir string
	// FileLevels can use multi file for record level logs. eg:
	//
	//  "error.log": []slog.Level{slog.Warn, slog.Error},
	//  "info.log": []slog.Level{slog.Trace, slog.Info, slog.Notice}
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

// Close handle
func (h *MultiFileHandler) Close() error {
	panic("implement me")
}

// Flush handle
func (h *MultiFileHandler) Flush() error {
	panic("implement me")
}

// Handle log record
func (h *MultiFileHandler) Handle(_ *slog.Record) error {
	panic("implement me")
}
