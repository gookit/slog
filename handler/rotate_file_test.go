package handler_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
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

	err = l.Flush()
	assert.NoError(t, err)

	// checkLogFileContents(t, fpath)
}

func TestNewRotateFileHandler(t *testing.T) {
	// by size
	logfile := "./testdata/both-rotate-file1.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	h, err := handler.NewRotateFileHandler(logfile, handler.EveryMinute, handler.WithMaxSize(128))

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
	}
	l.MustFlush()

	// by time
	logfile = "./testdata/both-rotate-file2.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))

	h, err = handler.NewRotateFileHandler(logfile, handler.EverySecond)
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l = slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("second ", i+1)
		time.Sleep(time.Second * 1)
	}
	err = l.FlushAll()
	assert.NoError(t, err)
}

func TestNewTimeRotateFileHandler_EveryDay(t *testing.T) {
	logfile := "./testdata/time-rotate_EveryDay.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))
	newfile := logfile + timex.Now().DateFormat(".YMD")
	assert.NoError(t, fsutil.DeleteIfFileExist(newfile))

	sec := -2
	// set current time is today 23:23:57
	testClock := func() time.Time {
		// dump.P(sec)
		return timex.Now().DayEnd().AddSeconds(sec).Time
	}

	// backup
	bckFn := rotatefile.DefaultTimeClockFn
	rotatefile.DefaultTimeClockFn = testClock
	defer func() {
		rotatefile.DefaultTimeClockFn = bckFn
	}()

	options := []handler.ConfigFn{
		handler.WithBuffSize(128),
	}

	h, err := handler.NewTimeRotateFileHandler(logfile, handler.EveryDay, options...)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true
	l.TimeClock = testClock

	for i := 0; i < 4; i++ {
		l.WithData(sampleData).Info("the th:", i, "info message")
		l.Warnf("the th:%d warn message text", i)
		sec++
		fmt.Println("log number ", (i+1)*2)
		// time.Sleep(time.Second * 1)
	}

	l.MustFlush()
	checkLogFileContents(t, logfile)
	checkLogFileContents(t, newfile)
}

func TestNewTimeRotateFileHandler_EveryHour(t *testing.T) {
	logfile := "./testdata/time-rotate_EveryHour.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(logfile))
	h, err := handler.NewTimeRotateFileHandler(logfile, handler.EveryHour)

	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("number ", i+1)
		// time.Sleep(time.Second * 1)
	}
	l.MustFlush()

	checkLogFileContents(t, logfile)
}

func checkLogFileContents(t *testing.T, logfile string) {
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
	l.MustFlush()
	// assert.NoError(t, os.Remove(fpath))
}
