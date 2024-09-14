package rotatefile_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/rotatefile"
)

// https://github.com/gookit/slog/issues/138
// 日志按everyday自动滚动，文件名的日期对应的是前一天的日志 #138
func TestIssues_138(t *testing.T) {
	logfile := "testdata/rotate_day.log"

	mt := rotatefile.NewMockClock("2023-11-16 23:59:55")
	w, err := rotatefile.NewWriterWith(rotatefile.WithDebugMode, func(c *rotatefile.Config) {
		c.TimeClock = mt
		// c.MaxSize = 128
		c.Filepath = logfile
		c.RotateTime = rotatefile.EveryDay
	})

	assert.NoErr(t, err)
	defer w.MustClose()

	for i := 0; i < 15; i++ {
		dt := mt.Datetime()
		_, err = w.WriteString(dt + " [INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
		// increase time
		mt.Add(time.Second * 3)
		// mt.Add(time.Millisecond * 300)
	}

	// Out: rotate_day.log, rotate_day.log.20231116
	files := fsutil.Glob(logfile + "*")
	assert.Len(t, files, 2)

	// check contents
	assert.True(t, fsutil.IsFile(logfile))
	s := fsutil.ReadString(logfile)
	assert.StrContains(t, s, "2023-11-17 00:00")

	oldFile := logfile + ".20231116"
	assert.True(t, fsutil.IsFile(oldFile))
	s = fsutil.ReadString(oldFile)
	assert.StrContains(t, s, "2023-11-16 23:")
}

// https://github.com/gookit/slog/issues/150
// 日志轮转时间设置为分钟时，FirstCheckTime计算单位错误，导致生成预期外的多个日志文件 #150
func TestIssues_150(t *testing.T) {
	logfile := "testdata/i150_rotate_min.log"

	mt := rotatefile.NewMockClock("2024-09-14 18:39:55")
	w, err := rotatefile.NewWriterWith(rotatefile.WithDebugMode, func(c *rotatefile.Config) {
		c.TimeClock = mt
		// c.MaxSize = 128
		c.Filepath = logfile
		c.RotateTime = rotatefile.EveryMinute * 3
	})

	assert.NoErr(t, err)
	defer w.MustClose()

	for i := 0; i < 15; i++ {
		dt := mt.Datetime()
		_, err = w.WriteString(dt + " [INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
		// increase time
		mt.Add(time.Minute * 1)
	}

	// Out: rotate_day.log, rotate_day.log.20231116
	files := fsutil.Glob(logfile + "*")
	assert.LenGt(t, files, 3)

	// check contents
	assert.True(t, fsutil.IsFile(logfile))
	s := fsutil.ReadString(logfile)
	assert.StrContains(t, s, "2024-09-14 18:")

	// i150_rotate_min.log.20240914_1842
	oldFile := logfile + ".20240914_1842"
	assert.True(t, fsutil.IsFile(oldFile))
	s = fsutil.ReadString(oldFile)
	assert.StrContains(t, s, "2024-09-14 18:41")
}
