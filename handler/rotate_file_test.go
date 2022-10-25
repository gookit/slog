package handler_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func TestNewRotateFileHandler(t *testing.T) {
	// by size
	logfile := "./testdata/both-rotate-bysize.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h, err := handler.NewRotateFileHandler(logfile, handler.EveryMinute, handler.WithMaxSize(128))

	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
	}
	l.MustFlush()

	// by time
	logfile = "./testdata/both-rotate-bytime.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h, err = handler.NewRotateFileHandler(logfile, handler.EverySecond)
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l = slog.NewWithHandlers(h)

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("second ", i+1)
		time.Sleep(time.Second * 1)
	}
	l.Error("error message")

	assert.NoErr(t, l.FlushAll())
}

func TestNewSizeRotateFileHandler(t *testing.T) {
	logfile := "./testdata/size-rotate-file.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h, err := handler.NewSizeRotateFileHandler(logfile, 468, handler.WithBuffSize(256))
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true
	l.CallerFlag = slog.CallerFlagFull

	for i := 0; i < 4; i++ {
		l.Info("this is a info", "message, index=", i)
		l.Warn("this is a warn message, index=", i)
	}

	assert.NoErr(t, l.Close())
	checkLogFileContents(t, logfile)
}

func TestNewTimeRotateFileHandler_EveryDay(t *testing.T) {
	logfile := "./testdata/time-rotate_EveryDay.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))
	newFile := logfile + timex.Now().DateFormat(".Ymd")
	assert.NoErr(t, fsutil.DeleteIfFileExist(newFile))

	sec := -2
	// set current time to today 23:59:57
	testClock := func() time.Time {
		// dump.P(sec)
		return timex.Now().DayEnd().AddSeconds(sec).Time
	}
	assert.Eq(t, "23:59:57", timex.Date(testClock(), "H:I:S"))

	// backup
	bckFn := rotatefile.DefaultTimeClockFn
	rotatefile.DefaultTimeClockFn = testClock
	defer func() {
		rotatefile.DefaultTimeClockFn = bckFn
	}()

	options := []handler.ConfigFn{
		handler.WithBuffSize(128),
	}

	h := handler.MustTimeRotateFile(logfile, handler.EveryDay, options...)
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
	checkLogFileContents(t, newFile)
}

func TestNewTimeRotateFileHandler_EveryHour(t *testing.T) {
	logfile := "./testdata/time-rotate_EveryHour.log"
	assert.NoErr(t, fsutil.DeleteIfExist(logfile))

	hourStart := timex.Now().HourStart()
	newFile := logfile + hourStart.DateFormat(".Ymd_H00")
	assert.NoErr(t, fsutil.DeleteIfFileExist(newFile))

	sec := -2
	// set current time to hour end 59:58
	testClock := func() time.Time {
		// dump.P(sec)
		return hourStart.AddHour(1).AddSeconds(sec).Time
	}
	assert.Eq(t, "59:58", timex.Date(testClock(), "I:S"))

	// backup
	bckFn := rotatefile.DefaultTimeClockFn
	rotatefile.DefaultTimeClockFn = testClock
	defer func() {
		rotatefile.DefaultTimeClockFn = bckFn
	}()

	options := []handler.ConfigFn{
		handler.WithBuffSize(0),
	}
	h, err := handler.NewTimeRotateFile(logfile, rotatefile.EveryHour, options...)

	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true
	l.TimeClock = testClock

	for i := 0; i < 3; i++ {
		l.WithData(sampleData).Info("the th:", i, "info message")
		l.Warnf("the th:%d warn message text", i)
		sec++
		fmt.Println("log number ", (i+1)*2)
	}
	l.MustFlush()

	checkLogFileContents(t, logfile)
	checkLogFileContents(t, newFile)
}

func TestNewTimeRotateFileHandler_someSeconds(t *testing.T) {
	logfile := "./testdata/time-rotate-Seconds.log"
	assert.NoErr(t, fsutil.DeleteIfExist(logfile))
	h, err := handler.NewTimeRotateFileHandler(logfile, handler.EverySecond)

	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
		fmt.Println("second ", i+1)
		time.Sleep(time.Second * 1)
	}
	l.MustFlush()
	// assert.NoErr(t, os.Remove(fpath))
}

func checkLogFileContents(t *testing.T, logfile string) {
	assert.True(t, fsutil.IsFile(logfile))

	bts, err := ioutil.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARNING]")
	assert.Contains(t, str, "warn message")
}
