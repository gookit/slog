package handler

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/rotatefile"
)

// MustRotateFile handler instance
func MustRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) *SyncCloseHandler {
	h, err := NewRotateFileHandler(logfile, rt, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewRotateFile instance
func NewRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error) {
	return NewRotateFileHandler(logfile, rt, fns...)
}

// NewRotateFileHandler instance
// It supports splitting log files by time and size
func NewRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error) {
	cfg := NewConfig(fns...).With(WithLogfile(logfile), WithRotateTime(rt))

	writer, err := cfg.RotateWriter()
	if err != nil {
		return nil, err
	}

	h := &SyncCloseHandler{
		Output: writer,
		// with log levels and formatter
		LevelFormattable: slog.NewLvsFormatter(cfg.Levels),
	}

	return h, nil
}

//
// ---------------------------------------------------------------------------
// rotate file by size
// ---------------------------------------------------------------------------
//

// MustSizeRotateFile instance
func MustSizeRotateFile(logfile string, size int, fns ...ConfigFn) *SyncCloseHandler {
	h, err := NewSizeRotateFileHandler(logfile, size, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewSizeRotateFile instance
func NewSizeRotateFile(logfile string, maxSize int, fns ...ConfigFn) (*SyncCloseHandler, error) {
	return NewSizeRotateFileHandler(logfile, maxSize, fns...)
}

// NewSizeRotateFileHandler instance
func NewSizeRotateFileHandler(logfile string, maxSize int, fns ...ConfigFn) (*SyncCloseHandler, error) {
	// close rotate by time.
	fns = append(fns, WithMaxSize(uint64(maxSize)))

	return NewRotateFileHandler(logfile, 0, fns...)
}

//
// ---------------------------------------------------------------------------
// rotate log file by time
// ---------------------------------------------------------------------------
//

// RotateTime rotate log file by time.
//
// EveryDay:
//   - "error.log.20201223"
//
// EveryHour, Every30Minutes, EveryMinute:
//   - "error.log.20201223_1500"
//   - "error.log.20201223_1530"
//   - "error.log.20201223_1523"
//
// Deprecated: please use rotatefile.RotateTime
type RotateTime = rotatefile.RotateTime

// Deprecated: Please use define constants on pkg rotatefile. e.g. rotatefile.EveryDay
const (
	EveryDay  = rotatefile.EveryDay
	EveryHour = rotatefile.EveryDay

	Every30Minutes = rotatefile.Every30Min
	Every15Minutes = rotatefile.Every15Min

	EveryMinute = rotatefile.EveryMinute
	EverySecond = rotatefile.EverySecond // only use for tests
)

// MustTimeRotateFile instance
func MustTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) *SyncCloseHandler {
	h, err := NewTimeRotateFileHandler(logfile, rt, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewTimeRotateFile instance
func NewTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error) {
	return NewTimeRotateFileHandler(logfile, rt, fns...)
}

// NewTimeRotateFileHandler instance
func NewTimeRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error) {
	// close rotate by size.
	fns = append(fns, WithMaxSize(0))

	return NewRotateFileHandler(logfile, rt, fns...)
}
