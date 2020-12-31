package handler

import (
	"os"
	"time"

	"github.com/gookit/slog"
)

// RotateFileHandler struct definition
// It also supports splitting log files by time and size
type RotateFileHandler struct {
	lockWrapper
	bufFileWrapper

	// LevelsWithFormatter support limit log levels and formatter
	LevelsWithFormatter

	// for size rotating
	written   uint64
	rotateNum uint

	// for time rotating
	rotateTime     rotateTime
	suffixFormat   string
	checkInterval  int64
	nextRotatingAt int64

	// for clear log files
	MaxFileCount int // The number of files should be kept
	MaxKeepTime  int // Time to wait until old logs are purged.

	// file contents max size
	MaxSize uint64
	// RenameFunc build filename for rotate file
	RenameFunc func(fpath string, rotateNum uint) string
}

// MustRotateFile instance
func MustRotateFile(filepath string, rt rotateTime) *RotateFileHandler {
	h, err := NewRotateFileHandler(filepath, rt)
	if err != nil {
		panic(err)
	}

	return h
}

// NewRotateFile instance
func NewRotateFile(filepath string, rt rotateTime) (*RotateFileHandler, error) {
	return NewRotateFileHandler(filepath, rt)
}

// NewRotateFileHandler instance
func NewRotateFileHandler(filepath string, rt rotateTime) (*RotateFileHandler, error) {
	h := &RotateFileHandler{
		rotateTime: rt,
		// file contents size
		MaxSize: DefaultMaxSize,
		// default log all levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
		// build new log filename.
		// eg: "error.log" => "error.log.010215_00001"
		RenameFunc: defaultNewLogfileFunc,
	}

	// init
	h.checkInterval, h.suffixFormat = rt.GetIntervalAndFormat()

	// fw := fileWrapper{fpath: filepath}
	fw := bufFileWrapper{
		BuffSize: defaultBufferSize,
	}
	fw.fpath = filepath
	// set prop
	h.bufFileWrapper = fw

	// open log file
	if err := h.ReopenFile(); err != nil {
		return nil, err
	}

	// storage next rotating time
	fileInfo, err := h.file.Stat()
	if err != nil {
		return nil, err
	}

	h.nextRotatingAt = fileInfo.ModTime().Unix() + h.checkInterval

	return h, nil
}

// func (h *RotateFileHandler) Write(p []byte) (n int, err error) {
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
func (h *RotateFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return err
	}

	// if lock enabled
	h.Lock()
	defer h.Unlock()

	var n int

	// write logs
	n, err = h.Write(bts)
	if err == nil {
		h.written += uint64(n)

		// do rotating file by time
		err = h.byTimeRotatingFile()

		// do rotating file by size
		if h.written >= h.MaxSize {
			err = h.bySizeRotatingFile()
		}
	}
	return
}

func (h *RotateFileHandler) byTimeRotatingFile() error {
	now := time.Now()
	if h.nextRotatingAt > now.Unix() {
		return nil
	}

	// close file
	if err := h.Close(); err != nil {
		return err
	}

	// rename current to new file
	newFilepath := h.fpath + "." + now.Format(h.suffixFormat)

	// do rotating file
	err := h.doRotatingFile(newFilepath)

	// storage next rotating time
	h.nextRotatingAt = now.Unix() + h.checkInterval
	return err
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (h *RotateFileHandler) bySizeRotatingFile() error {
	// close file
	if err := h.Close(); err != nil {
		return err
	}

	// rename current to new file
	h.rotateNum++
	newFilepath := h.RenameFunc(h.fpath, h.rotateNum)

	// do rotating file
	return h.doRotatingFile(newFilepath)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (h *RotateFileHandler) doRotatingFile(newFilepath string) error {
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
