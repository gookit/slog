package handler

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gookit/slog"
)

// default new Logfile Func
var defaultNewLogfileFunc = func(fpath string, rotateNum uint) string {
	// 2006-01-02 15-04-05
	suffix := time.Now().Format("010215")
	return fmt.Sprintf("%s.%s_%05d", fpath, suffix, rotateNum)
}

// SizeRotateFileHandler struct definition
type SizeRotateFileHandler struct {
	FileHandler

	written   uint64
	rotateNum uint
	// file contents max size
	MaxSize uint64
	// RenameFunc build filename for rotate file
	RenameFunc func(fpath string, rotateNum uint) string
}

// MustSizeRotateFile instance
func MustSizeRotateFile(logfile string, maxSize uint64) *SizeRotateFileHandler {
	h, err := NewSizeRotateFileHandler(logfile, maxSize)
	if err != nil {
		panic(err)
	}

	return h
}

// NewSizeRotateFile instance
func NewSizeRotateFile(logfile string, maxSize uint64) (*SizeRotateFileHandler, error) {
	return NewSizeRotateFileHandler(logfile, maxSize)
}

// NewSizeRotateFileHandler instance
func NewSizeRotateFileHandler(logfile string, maxSize uint64) (*SizeRotateFileHandler, error) {
	h := &SizeRotateFileHandler{
		// MaxSize: DefaultMaxSize,
		MaxSize: maxSize,
		// init file handler
		FileHandler: FileHandler{
			fpath: logfile,
			// buffer size
			BuffSize: defaultBufferSize,
			// default log all levels
			LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
		},
		// build new filename.
		// eg: "error.log" => "error.log.010215_00001"
		RenameFunc: defaultNewLogfileFunc,
	}

	file, err := QuickOpenFile(logfile)
	if err != nil {
		return nil, err
	}

	h.file = file
	return h, nil
}

// func (h *SizeRotateFileHandler) Write(p []byte) (n int, err error) {
// 	if h.written+uint64(len(p)) >= h.MaxSize {
// 		if err := h.rotateFile(time.Now()); err != nil {
// 			return 0, err
// 		}
// 	}
//
// 	n, err = h.file.Write(p)
// 	h.written += uint64(n)
// 	return
// }

// Handle the log record
func (h *SizeRotateFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return err
	}

	// if lock enabled
	h.Lock()
	defer h.Unlock()

	var n int

	// direct write logs to file
	if h.NoBuffer {
		n, err = h.file.Write(bts)
	} else {
		// enable buffer
		if h.bufio == nil {
			h.bufio = bufio.NewWriterSize(h.file, h.BuffSize)
		}

		n, err = h.bufio.Write(bts)
	}

	if err == nil {
		h.written += uint64(n)

		// do rotating file
		if h.written >= h.MaxSize {
			err = h.bySizeRotatingFile()
		}
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (h *SizeRotateFileHandler) bySizeRotatingFile() error {
	// close file
	if err := h.Close(); err != nil {
		return err
	}

	// rename current to new file
	h.rotateNum++
	newFilepath := h.RenameFunc(h.fpath, h.rotateNum)
	err := os.Rename(h.fpath, newFilepath)
	if err != nil {
		return err
	}

	// reopen file
	h.file, err = QuickOpenFile(h.fpath)
	if err != nil {
		return err
	}

	// if enable buffer
	if h.bufio != nil {
		h.bufio.Reset(h.file)
	}

	// reset h.written
	h.written = 0
	return nil
}
