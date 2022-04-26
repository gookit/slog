package handler

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
)

// MustSimpleFile new instance
func MustSimpleFile(filepath string) *SyncCloseHandler {
	h, err := NewSimpleFileHandler(filepath)
	if err != nil {
		panic(err)
	}

	return h
}

// NewSimpleFile new instance
func NewSimpleFile(filepath string) (*SyncCloseHandler, error) {
	return NewSimpleFileHandler(filepath)
}

// NewSimpleFileHandler instance
//
// Usage:
// 	h, err := NewSimpleFileHandler("/tmp/error.log")
//
// custom formatter
//	h.SetFormatter(slog.NewJSONFormatter())
//	slog.PushHandler(h)
//	slog.Info("log message")
func NewSimpleFileHandler(filePath string) (*SyncCloseHandler, error) {
	file, err := fsutil.QuickOpenFile(filePath)
	if err != nil {
		return nil, err
	}

	h := &SyncCloseHandler{
		Output: file,
		// init default log level
		LevelFormattable: slog.NewLvFormatter(slog.InfoLevel),
	}

	return h, nil
}
