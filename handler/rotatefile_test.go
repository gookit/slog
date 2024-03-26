package handler_test

import (
	"fmt"
	"os"
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

	h, err := handler.NewRotateFile(logfile, handler.EveryMinute, handler.WithMaxSize(128))
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true

	for i := 0; i < 3; i++ {
		l.Info("info", "message", i)
		l.Warn("warn message", i)
	}
	l.MustClose()

	// by time
	logfile = "./testdata/both-rotate-bytime.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

	h = handler.MustRotateFile(logfile, handler.EverySecond)
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
	t.Run("NewSizeRotateFile", func(t *testing.T) {
		logfile := "./testdata/size-rotate-file.log"
		assert.NoErr(t, fsutil.DeleteIfFileExist(logfile))

		h, err := handler.NewSizeRotateFile(logfile, 468, handler.WithBuffSize(256))
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
	})

	t.Run("MustSizeRotateFile", func(t *testing.T) {
		logfile := "./testdata/must-size-rotate-file.log"
		h := handler.MustSizeRotateFile(logfile, 128, handler.WithBuffSize(128))
		h.SetFormatter(slog.NewJSONFormatter())
		err := h.Handle(newLogRecord("this is a info message"))
		assert.NoErr(t, err)

		files := fsutil.Glob(logfile + "*")
		assert.Len(t, files, 2)
	})
}

func TestNewTimeRotateFileHandler_EveryDay(t *testing.T) {
	logfile := "./testdata/time-rotate_EveryDay.log"
	newFile := logfile + ".20221116"

	clock := rotatefile.NewMockClock("2022-11-16 23:59:57")
	options := []handler.ConfigFn{
		handler.WithBuffSize(128),
		handler.WithTimeClock(clock),
	}

	h := handler.MustTimeRotateFile(logfile, handler.EveryDay, options...)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true
	l.TimeClock = clock.Now

	for i := 0; i < 6; i++ {
		l.WithData(sampleData).Info("the th:", i, "info message")
		l.Warnf("the th:%d warn message text", i)
		fmt.Println("log number ", (i+1)*2)
		clock.Add(time.Second * 1)
	}

	l.MustClose()
	checkLogFileContents(t, logfile)
	checkLogFileContents(t, newFile)
}

func TestNewTimeRotateFileHandler_EveryHour(t *testing.T) {
	clock := rotatefile.NewMockClock("2022-04-28 20:59:58")
	logfile := "./testdata/time-rotate_EveryHour.log"
	newFile := logfile + timex.DateFormat(clock.Now(), ".Ymd_H00")

	options := []handler.ConfigFn{
		handler.WithTimeClock(clock),
		handler.WithBuffSize(0),
	}
	h, err := handler.NewTimeRotateFile(logfile, rotatefile.EveryHour, options...)

	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(logfile))

	l := slog.NewWithHandlers(h)
	l.ReportCaller = true
	l.TimeClock = clock.Now

	for i := 0; i < 6; i++ {
		l.WithData(sampleData).Info("the th:", i, "info message")
		l.Warnf("the th:%d warn message text", i)
		fmt.Println("log number ", (i+1)*2)
		clock.Add(time.Second * 1)
	}
	l.MustClose()

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
	l.MustClose()
	// assert.NoErr(t, os.Remove(fpath))
}

func checkLogFileContents(t *testing.T, logfile string) {
	assert.True(t, fsutil.IsFile(logfile))

	bts, err := os.ReadFile(logfile)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARN]")
	assert.Contains(t, str, "warn message")
}
