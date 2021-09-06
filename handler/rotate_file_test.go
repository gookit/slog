package handler_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewSizeRotateFileHandler(t *testing.T) {
	fpath := "./testdata/size-rotate-file.log"
	h, err := handler.NewSizeRotateFileHandler(fpath, 128)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
	}
	l.Flush()
	// checkLogFileContents(t, fpath)
}

func TestNewRotateFileHandler(t *testing.T) {
	// by size
	fpath := "./testdata/both-rotate-file1.log"
	h, err := handler.NewRotateFileHandler(fpath, handler.EveryMinute)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	h.MaxSize = 128

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
	}
	l.Flush()

	// by time
	fpath = "./testdata/both-rotate-file2.log"
	h, err = handler.NewRotateFileHandler(fpath, handler.EverySecond)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	l = slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("second ", i+1)
		time.Sleep(time.Second * 1)
	}
	l.Flush()
}

func TestNewTimeRotateFileHandler_EveryDay(t *testing.T) {
	fpath := "./testdata/time-rotate-file_EveryDay.log"
	h, err := handler.NewTimeRotateFileHandler(fpath, handler.EveryDay)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("number ", i+1)
		// time.Sleep(time.Second * 1)
	}
	l.Flush()

	checkLogFileContents(t, fpath)
}

func TestNewTimeRotateFileHandler_EveryHour(t *testing.T) {
	fpath := "./testdata/time-rotate-file_EveryHour.log"
	h, err := handler.NewTimeRotateFileHandler(fpath, handler.EveryHour)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("number ", i+1)
		// time.Sleep(time.Second * 1)
	}
	l.Flush()

	checkLogFileContents(t, fpath)
}

func checkLogFileContents(t *testing.T, logfile string)  {
	assert.True(t, fsutil.IsFile(logfile))

	bts, err := ioutil.ReadFile(logfile)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARNING]")
	assert.Contains(t, str, "warn message")
}

func TestNewTimeRotateFileHandler(t *testing.T) {
	fpath := "./testdata/time-rotate-file.log"
	h, err := handler.NewTimeRotateFileHandler(fpath, handler.EverySecond)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(fpath))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("second ", i+1)
		time.Sleep(time.Second * 1)
	}
	l.Flush()
	// assert.NoError(t, os.Remove(fpath))
}
