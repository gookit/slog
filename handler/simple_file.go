package handler

import (
	"github.com/gookit/slog"
)

// SimpleFileHandler struct
type SimpleFileHandler struct {
	fileHandler
	lockWrapper
	// LevelWithFormatter support level and formatter
	LevelWithFormatter
}

// NewSimpleFileHandler instance
//
// Usage:
// 	h, err := NewSimpleFileHandler("", DefaultFileFlags)
// custom file flags
// 	h, err := NewSimpleFileHandler("", os.O_CREATE | os.O_WRONLY | os.O_APPEND)
// custom formatter
//	h.SetFormatter(slog.NewJSONFormatter())
//	slog.PushHandler(h)
//	slog.Info("log message")
func NewSimpleFileHandler(filepath string) (*SimpleFileHandler, error) {
	fh := fileHandler{fpath: filepath}
	if err := fh.ReopenFile(); err != nil {
		return nil, err
	}

	h := &SimpleFileHandler{
		fileHandler: fh,
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

	// if enable lock
	if h.LockEnabled() {
		h.Lock()
		defer h.Unlock()
	}

	// direct write logs
	_, err = h.file.Write(bts)
	return
}

// Close handler, will be flush logs to file, then close file
func (h *SimpleFileHandler) Close() error {
	return h.fileHandler.Close()
}

// Flush logs to disk file
func (h *SimpleFileHandler) Flush() error {
	return h.fileHandler.Flush()
}
