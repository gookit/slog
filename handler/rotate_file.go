package handler

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gookit/slog"
)

// RotateFileHandler struct definition
type RotateFileHandler struct {
	FileHandler
	baseFile  string
	rotateNum uint

	written uint64
	// file contents max size
	MaxSize uint64
	// RenameFunc build filename for rotate file
	RenameFunc func(baseFile string, rotateNum uint) string
}

// NewRotateFileHandler instance
func NewRotateFileHandler(filepath string, maxSize uint64) (*RotateFileHandler, error) {
	h := &RotateFileHandler{
		// MaxSize: DefaultMaxSize,
		MaxSize:  maxSize,
		baseFile: filepath,
		FileHandler: FileHandler{
			fpath: filepath,
			// buffer size
			BuffSize: defaultBufferSize,
			// init log levels
			LevelsWithFormatter: LevelsWithFormatter{
				Levels: slog.AllLevels, // default log all levels
			},
		},
		// build new filename. eg: "error.log" => "error.log.0102150405"
		RenameFunc: func(baseFile string, rotateNum uint) string {
			// 2006-01-02 15-04-05
			// suffix := time.Now().Format("0102150405")
			return fmt.Sprintf("%s.%04d", baseFile, rotateNum)
		},
	}

	file, err := openFile(filepath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return nil, err
	}

	h.file = file
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
			err = h.doRotatingFile()
		}
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (h *RotateFileHandler) doRotatingFile() error {
	// if h.written < h.MaxSize {
	// 	return nil
	// }

	// close file
	if err := h.Close(); err != nil {
		return err
	}

	// rename current to new file
	h.rotateNum++
	newFilepath := h.RenameFunc(h.baseFile, h.rotateNum)
	err := os.Rename(h.baseFile, newFilepath)
	if err != nil {
		return err
	}

	// reopen file
	h.file, err = openFile(h.baseFile, DefaultFileFlags, DefaultFilePerm)
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
	EverySecond // only use for tests
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
	case EverySecond:
		return "Every Second"
	}

	// should never return this
	return "Unknown"
}

func (rt rotateType) GetIntervalAndFormat() (checkInterval int64, suffixFormat string) {
	switch rt {
	case EveryDay:
		checkInterval = 3600 * 24
		suffixFormat = "20060102"
	case EveryHour:
		checkInterval = 3600
		suffixFormat = "20060102_1500"
	case Every30Minutes:
		checkInterval = 1800
		suffixFormat = "20060102_1504"
	case EveryMinute:
		checkInterval = 60
		suffixFormat = "20060102_1504"
	}

	// Every Second
	return 1, "20060102_150405"
}

// TimeRotateFileHandler struct
// refer http://hg.python.org/cpython/file/2.7/Lib/logging/handlers.py
// refer https://github.com/flike/golog/blob/master/filehandler.go
type TimeRotateFileHandler struct {
	lockWrapper
	bufFileWrapper

	// LevelsWithFormatter support limit log levels and formatter
	LevelsWithFormatter

	// file *os.File
	// baseFile string

	rotateType   rotateType
	suffixFormat string

	checkInterval  int64
	nextRotatingAt int64
}

// NewTimeRotateFileHandler instance
func NewTimeRotateFileHandler(filepath string, rt rotateType) (*TimeRotateFileHandler, error) {
	h := &TimeRotateFileHandler{
		rotateType: rt,
		// init log levels
		LevelsWithFormatter: LevelsWithFormatter{
			Levels: slog.AllLevels, // default log all levels
		},
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

// Handle the log record
func (h *TimeRotateFileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// if enable lock
	h.Lock()
	defer h.Unlock()

	// direct write logs to file
	if h.NoBuffer {
		_, err = h.file.Write(bts)
	} else {
		// enable buffer
		if h.bufio == nil {
			h.bufio = bufio.NewWriterSize(h.file, h.BuffSize)
		}

		_, err = h.bufio.Write(bts)
	}

	// do rotating file
	if err == nil {
		err = h.doRotatingFile()
	}
	return
}

func (h *TimeRotateFileHandler) doRotatingFile() error {
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
	err := os.Rename(h.fpath, newFilepath)
	if err != nil {
		return err
	}

	// reopen file
	h.file, err = openFile(h.fpath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return err
	}

	// if enable buffer
	if h.bufio != nil {
		h.bufio.Reset(h.file)
	}

	// storage next rotating time
	h.nextRotatingAt = time.Now().Unix() + h.checkInterval
	return nil
}
