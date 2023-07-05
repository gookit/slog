package handler

import (
	"github.com/gookit/goutil/basefn"
	"github.com/gookit/slog"
)

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(logfile string, fns ...ConfigFn) (*SyncCloseHandler, error) {
	return NewFileHandler(logfile, append(fns, WithUseJSON(true))...)
}

// NewBuffFileHandler create file handler with buff size
func NewBuffFileHandler(logfile string, buffSize int, fns ...ConfigFn) (*SyncCloseHandler, error) {
	return NewFileHandler(logfile, append(fns, WithBuffSize(buffSize))...)
}

// MustFileHandler create file handler
func MustFileHandler(logfile string, fns ...ConfigFn) *SyncCloseHandler {
	return basefn.Must(NewFileHandler(logfile, fns...))
}

// NewFileHandler create new FileHandler
func NewFileHandler(logfile string, fns ...ConfigFn) (h *SyncCloseHandler, err error) {
	return NewEmptyConfig(fns...).With(WithLogfile(logfile)).CreateHandler()
}

//
// ------------- simple file handler -------------
//

// MustSimpleFile new instance
func MustSimpleFile(filepath string, maxLv ...slog.Level) *SyncCloseHandler {
	return basefn.Must(NewSimpleFileHandler(filepath, maxLv...))
}

// NewSimpleFile new instance
func NewSimpleFile(filepath string, maxLv ...slog.Level) (*SyncCloseHandler, error) {
	return NewSimpleFileHandler(filepath, maxLv...)
}

// NewSimpleFileHandler instance, default log level is InfoLevel
//
// Usage:
//
//	h, err := NewSimpleFileHandler("/tmp/error.log")
//
// Custom formatter:
//
//	h.SetFormatter(slog.NewJSONFormatter())
//	slog.PushHandler(h)
//	slog.Info("log message")
func NewSimpleFileHandler(filePath string, maxLv ...slog.Level) (*SyncCloseHandler, error) {
	file, err := QuickOpenFile(filePath)
	if err != nil {
		return nil, err
	}

	h := SyncCloserWithMaxLevel(file, basefn.FirstOr(maxLv, slog.InfoLevel))
	return h, nil
}
