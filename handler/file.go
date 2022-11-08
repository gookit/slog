package handler

import (
	"github.com/gookit/slog"
)

// MustFileHandler create file handler
func MustFileHandler(logfile string, fns ...ConfigFn) *SyncCloseHandler {
	h, err := NewFileHandler(logfile, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(logfile string, fns ...ConfigFn) (*SyncCloseHandler, error) {
	fns = append(fns, WithUseJSON(true))

	return NewFileHandler(logfile, fns...)
}

// NewBuffFileHandler create file handler with buff size
func NewBuffFileHandler(logfile string, buffSize int, fns ...ConfigFn) (*SyncCloseHandler, error) {
	fns = append(fns, WithBuffSize(buffSize))

	return NewFileHandler(logfile, fns...)
}

// NewFileHandler create new FileHandler
func NewFileHandler(logfile string, fns ...ConfigFn) (h *SyncCloseHandler, err error) {
	cfg := NewEmptyConfig(fns...).With(WithLogfile(logfile))

	return cfg.CreateHandler()
}

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
//
//	h, err := NewSimpleFileHandler("/tmp/error.log")
//
// custom formatter
//
//	h.SetFormatter(slog.NewJSONFormatter())
//	slog.PushHandler(h)
//	slog.Info("log message")
func NewSimpleFileHandler(filePath string) (*SyncCloseHandler, error) {
	file, err := QuickOpenFile(filePath)
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
