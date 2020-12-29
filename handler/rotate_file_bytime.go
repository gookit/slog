package handler

import (
	"os"
	"time"

	"github.com/gookit/slog"
)

// rotate log file by time.
//
// EveryDay:
// 	- "error.log.20201223"
// EveryHour, Every30Minutes,EveryMinute:
// 	- "error.log.20201223_1500"
// 	- "error.log.20201223_1530"
// 	- "error.log.20201223_1523"
type rotateTime uint8

const (
	EveryDay rotateTime = iota
	EveryHour
	Every30Minutes
	Every15Minutes
	EveryMinute
	EverySecond // only use for tests
)

// String rotate type to string
func (rt rotateTime) String() string {
	switch rt {
	case EveryDay:
		return "Every Day"
	case EveryHour:
		return "Every Hour"
	case Every30Minutes:
		return "Every 30 Minutes"
	case Every15Minutes:
		return "Every 15 Minutes"
	case EveryMinute:
		return "Every Minute"
	case EverySecond:
		return "Every Second"
	}

	// should never return this
	return "Unknown"
}

// GetIntervalAndFormat get check interval time and log suffix format
func (rt rotateTime) GetIntervalAndFormat() (checkInterval int64, suffixFormat string) {
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
	case Every15Minutes:
		checkInterval = 900
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

	rotateType   rotateTime
	suffixFormat string

	checkInterval  int64
	nextRotatingAt int64
}

// MustTimeRotateFile instance
func MustTimeRotateFile(filepath string, rt rotateTime) *TimeRotateFileHandler {
	h, err := NewTimeRotateFileHandler(filepath, rt)
	if err != nil {
		panic(err)
	}

	return h
}

// NewTimeRotateFile instance
func NewTimeRotateFile(filepath string, rt rotateTime) (*TimeRotateFileHandler, error) {
	return NewTimeRotateFileHandler(filepath, rt)
}

// NewTimeRotateFileHandler instance
func NewTimeRotateFileHandler(filepath string, rt rotateTime) (*TimeRotateFileHandler, error) {
	h := &TimeRotateFileHandler{
		rotateType: rt,
		// default log all levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
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

	// write logs
	_, err = h.Write(bts)

	// do rotating file
	if err == nil {
		err = h.byTimeRotatingFile()
	}
	return
}

func (h *TimeRotateFileHandler) byTimeRotatingFile() error {
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
	h.file, err = QuickOpenFile(h.fpath)
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
