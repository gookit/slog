package handler

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
)

// SimpleFileHandler struct
// - no buffer, will direct write logs to file.
type SimpleFileHandler struct {
	out SyncCloseWrapper
	LevelWithFormatter
}

// MustSimpleFile new instance
func MustSimpleFile(filepath string) *SimpleFileHandler {
	h, err := NewSimpleFileHandler(filepath)
	if err != nil {
		panic(err)
	}

	return h
}

// NewSimpleFile new instance
func NewSimpleFile(filepath string) (*SimpleFileHandler, error) {
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
func NewSimpleFileHandler(filePath string) (*SimpleFileHandler, error) {
	file, err := fsutil.QuickOpenFile(filePath)
	if err != nil {
		return nil, err
	}

	h := &SimpleFileHandler{
		out: NewSyncCloseWrapper(file),
		// init default log level
		LevelWithFormatter: newLvFormatter(slog.InfoLevel),
	}

	return h, nil
}

// Handle the log record
func (h *SimpleFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// direct write logs
	_, err = h.out.Write(bts)
	return
}
