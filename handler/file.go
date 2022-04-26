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
	cfg := NewConfig(fns...).With(WithLogfile(logfile))

	output, err := cfg.SyncCloseWriter()
	if err != nil {
		return nil, err
	}

	h = &SyncCloseHandler{
		Output: output,
		// with log levels and formatter
		LevelFormattable: slog.NewLvsFormatter(cfg.Levels),
	}

	if cfg.UseJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	}
	return h, nil
}
