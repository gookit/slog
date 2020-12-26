package handler

import (
	"io"
	"os"
	"path"

	"github.com/gookit/slog"
)

// SimpleFileHandler struct
type SimpleFileHandler struct {
	lockWrapper
	// LevelWithFormatter support level and formatter
	LevelWithFormatter
	// log file path. eg: "/var/log/my-app.log"
	fpath string
	file  *os.File
}

// NewSimpleFileHandler instance
// Usage:
// 	h, err := NewSimpleFileHandler("", DefaultFileFlags)
// custom file flags
// 	h, err := NewSimpleFileHandler("", os.O_CREATE | os.O_WRONLY | os.O_APPEND)
func NewSimpleFileHandler(filepath string, flag int) (*SimpleFileHandler, error) {
	dir := path.Dir(filepath)
	// if err := os.MkdirAll(dir, 0777); err != nil {
	if err := os.Mkdir(dir, 0777); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, flag, 0)
	if err != nil {
		return nil, err
	}

	h := &SimpleFileHandler{
		fpath: filepath,
		file:  file,
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

	// if use lock
	if h.LockEnabled() {
		h.Lock()
		defer h.Unlock()
	}

	// direct write logs
	_, err = h.file.Write(bts)
	return
}

// Writer return *os.File
func (h *SimpleFileHandler) Writer() io.Writer {
	return h.file
}

// Close handler, will be flush logs to file, then close file
func (h *SimpleFileHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.file.Close()
}

// Flush logs to disk file
func (h *SimpleFileHandler) Flush() error {
	return h.file.Sync()
}
