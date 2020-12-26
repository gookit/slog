package handler

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/gookit/slog"
)

// SizeRotateFileHandler struct
type SizeRotateFileHandler struct {
	// file contents max size
	MaxSize uint64
}

// log file examples:
// EveryDay:
// 	- "error.log.20201223"
// EveryHour, Every30Minutes,EveryMinute:
// 	- "error.log.20201223_1500"
// 	- "error.log.20201223_1530"
// 	- "error.log.20201223_1523"
type rotateType uint8
const (
	EveryDay rotateType = iota
	EveryHour
	Every30Minutes
	EveryMinute
)

// String rotate type to string
func (rt rotateType) String() string {
	switch rt {
	case EveryDay:
		return "Every Day"
	case EveryHour:
		return "Every Hour"
	case Every30Minutes:
		return "Every 30 Minutes"
	case EveryMinute:
		return "Every Minute"
	}

	// should never return this
	return "Unknown"
}

// TimeRotateFileHandler struct
// refer http://hg.python.org/cpython/file/2.7/Lib/logging/handlers.py
// refer https://github.com/flike/golog/blob/master/filehandler.go
type TimeRotateFileHandler struct {
	lockWrapper
	baseFile string

	rotateType rotateType
	suffixFormat string

	checkInterval  int64
	nextRotatingAt int64

	file *os.File
}

// NewTimeRotateFileHandler instance
func NewTimeRotateFileHandler(filepath string, rotateTime rotateType) (*TimeRotateFileHandler,error) {
	dir := path.Dir(filepath)
	// if err := os.MkdirAll(dir, 0777); err != nil {
	if err := os.Mkdir(dir, 0777); err != nil {
		return nil, err
	}

	h := &TimeRotateFileHandler{
		baseFile: filepath,
		rotateType: rotateTime,
	}

	switch rotateTime {
	case EveryDay:
		h.checkInterval = 3600 * 24
		h.suffixFormat = "20060102"
	case EveryHour:
		h.checkInterval = 3600
		h.suffixFormat = "20060102_1500"
	case Every30Minutes:
		h.checkInterval = 1800
		h.suffixFormat = "20060102_1504"
	case EveryMinute:
		h.checkInterval = 60
		h.suffixFormat = "20060102_1504"
	}

	file, err := os.OpenFile(filepath, DefaultFileFlags, 0666)
	if err != nil {
		return nil, err
	}

	h.file = file

	fInfo, _ := h.file.Stat()

	// storage next rotating time
	h.nextRotatingAt = fInfo.ModTime().Unix() + h.checkInterval

	return h, nil
}

func (h *TimeRotateFileHandler) doRotatingFile() error {
	now := time.Now()

	if h.nextRotatingAt <= now.Unix() {
		err := h.file.Close()
		if err != nil {
			return err
		}

		// move current to new file
		newFilepath := h.baseFile + "." + now.Format(h.suffixFormat)
		err = os.Rename(h.baseFile, newFilepath)
		if err != nil {
			return err
		}

		// reopen file
		h.file, err = os.OpenFile(h.baseFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		// storage next rotating time
		h.nextRotatingAt = time.Now().Unix() + h.checkInterval
	}

	return nil
}

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

// RotateFileHandler definition
type RotateFileHandler struct {
	FileHandler
	logger  *slog.Logger
	written uint64
	// RenameFunc for rotate file
	RenameFunc func(fpath string) string
}

func (h *RotateFileHandler) Write(p []byte) (n int, err error) {
	if h.written+uint64(len(p)) >= h.MaxSize {
		if err := h.rotateFile(time.Now()); err != nil {
			return 0, err
		}
	}

	n, err = h.file.Write(p)
	h.written += uint64(n)
	// if err != nil {
	// 	h.logger.Exit(err)
	// }
	return
}

// -------- refer from glog package
// rotateFile closes the syncBuffer's file and starts a new one.
func (h *RotateFileHandler) rotateFile(now time.Time) error {
	if h.file != nil {
		h.Flush()
		h.file.Close()
	}

	var err error
	h.file, _, err = create("INFO", now)
	h.written = 0
	if err != nil {
		return err
	}

	// init writer
	// h.Writer = bufio.NewWriterSize(h.file, bufferSize)

	// Write header.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
	n, err := h.file.Write(buf.Bytes())

	h.written += uint64(n)
	return err
}
