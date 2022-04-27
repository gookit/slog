package handler

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/rotatefile"
)

// RotateFileHandler struct definition
// It also supports splitting log files by time and size
type RotateFileHandler struct {
	// LockWrapper
	slog.LevelFormattable
	Output FlushCloseWriter
}

// MustRotateFile instance
func MustRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) *RotateFileHandler {
	h, err := NewRotateFileHandler(logfile, rt, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewRotateFile instance
func NewRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*RotateFileHandler, error) {
	return NewRotateFileHandler(logfile, rt, fns...)
}

// NewRotateFileHandler instance
func NewRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*RotateFileHandler, error) {
	cfg := NewConfig(fns...).With(WithLogfile(logfile), WithRotateTime(rt))

	writer, err := cfg.RotateWriter()
	if err != nil {
		return nil, err
	}

	h := &RotateFileHandler{
		Output: writer,
		// with log levels and formatter
		LevelFormattable: slog.NewLvsFormatter(cfg.Levels),
	}

	return h, nil
}

// Handle the log record
func (h *RotateFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return err
	}

	// if lock enabled
	// h.Lock()
	// defer h.Unlock()

	// write logs
	_, err = h.Output.Write(bts)
	return
}

// Close handler, will be flush logs to file, then close file
func (h *RotateFileHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.Output.Close()
}

// Flush logs to disk file
func (h *RotateFileHandler) Flush() error {
	return h.Output.Flush()
}

//
// ---------------------------------------------------------------------------
// rotate file by size
// ---------------------------------------------------------------------------
//

// MustSizeRotateFile instance
func MustSizeRotateFile(logfile string, size int, fns ...ConfigFn) *RotateFileHandler {
	h, err := NewSizeRotateFileHandler(logfile, size, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewSizeRotateFile instance
func NewSizeRotateFile(logfile string, maxSize int, fns ...ConfigFn) (*RotateFileHandler, error) {
	return NewSizeRotateFileHandler(logfile, maxSize, fns...)
}

// NewSizeRotateFileHandler instance
func NewSizeRotateFileHandler(logfile string, maxSize int, fns ...ConfigFn) (*RotateFileHandler, error) {
	// close rotate by time.
	fns = append(fns, WithMaxSize(maxSize))
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
// 	- "error.log.20201223"
// EveryHour, Every30Minutes, EveryMinute:
// 	- "error.log.20201223_1500"
// 	- "error.log.20201223_1530"
// 	- "error.log.20201223_1523"
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
func MustTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) *RotateFileHandler {
	h, err := NewTimeRotateFileHandler(logfile, rt, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewTimeRotateFile instance
func NewTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*RotateFileHandler, error) {
	return NewTimeRotateFileHandler(logfile, rt, fns...)
}

// NewTimeRotateFileHandler instance
func NewTimeRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*RotateFileHandler, error) {
	// close rotate by size.
	fns = append(fns, WithMaxSize(0))
	return NewRotateFileHandler(logfile, rt, fns...)
}
