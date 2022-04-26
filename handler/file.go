package handler

import (
	"io"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
)

// FileHandler definition
type FileHandler struct {
	// LevelsWithFormatter support limit log levels and formatter
	LevelsWithFormatter
	output SyncCloseWriter
}

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(logfile string) (*FileHandler, error) {
	return NewFileHandler(logfile, WithUseJSON(true))
}

// MustFileHandler create file handler
func MustFileHandler(logfile string, fns ...ConfigFn) *FileHandler {
	h, err := NewFileHandler(logfile, fns...)
	if err != nil {
		panic(err)
	}
	return h
}

// NewBuffFileHandler create file handler with buff size
func NewBuffFileHandler(filePath string, buffSize int, fns ...ConfigFn) (*FileHandler, error) {
	fns = append(fns, WithBuffSize(buffSize))

	return NewFileHandler(filePath, fns...)
}

// NewFileHandler create new FileHandler
func NewFileHandler(logfile string, fns ...ConfigFn) (h *FileHandler, err error) {
	cfg := NewConfig(fns...)
	cfg.Logfile = logfile

	var output SyncCloseWriter
	output, err = fsutil.QuickOpenFile(logfile)
	if err != nil {
		return nil, err
	}

	// wrap buffer
	if cfg.BuffSize > 0 {
		output = bufwrite.NewBufIOWriterSize(output, cfg.BuffSize)
	}

	h = &FileHandler{
		output: output,
		// with log levels and formatter
		LevelsWithFormatter: newLvsFormatter(cfg.Levels),
	}

	if cfg.UseJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	}

	return h, nil
}

// Writer return output writer
func (h *FileHandler) Writer() io.Writer {
	return h.output
}

// Close handler, will be flush logs to file, then close file
func (h *FileHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.output.Close()
}

// Flush logs to disk file
func (h *FileHandler) Flush() error {
	return h.output.Sync()
}

// Handle the log record
func (h *FileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// if enable lock
	// h.Lock()
	// defer h.Unlock()

	_, err = h.output.Write(bts)
	return
}
