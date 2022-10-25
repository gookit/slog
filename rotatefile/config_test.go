package rotatefile_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog/rotatefile"
)

func TestNewDefaultConfig(t *testing.T) {
	size := fmtutil.DataSize(1024 * 1024 * 10)
	dump.P(size)

	c := rotatefile.NewDefaultConfig()
	assert.Eq(t, rotatefile.DefaultMaxSize, c.MaxSize)
}

func TestNewConfig(t *testing.T) {
	cfg := rotatefile.NewConfig("testdata/test.log")

	assert.Eq(t, rotatefile.DefaultBackNum, cfg.BackupNum)
	assert.Eq(t, rotatefile.DefaultBackTime, cfg.BackupTime)
	assert.Eq(t, rotatefile.EveryHour, cfg.RotateTime)
	assert.Eq(t, rotatefile.DefaultMaxSize, cfg.MaxSize)

	dump.P(cfg)

	cfg = rotatefile.EmptyConfigWith(func(c *rotatefile.Config) {
		c.Compress = true
	})
	assert.True(t, cfg.Compress)
	assert.Eq(t, uint(0), cfg.BackupNum)
	assert.Eq(t, uint(0), cfg.BackupTime)
}

func TestRotateTime_String(t *testing.T) {
	assert.Eq(t, "Every 1 Day", rotatefile.EveryDay.String())
	assert.Eq(t, "Every 1 Hours", rotatefile.EveryHour.String())
	assert.Eq(t, "Every 1 Minutes", rotatefile.EveryMinute.String())
	assert.Eq(t, "Every 1 Seconds", rotatefile.EverySecond.String())

	assert.Eq(t, "Every 2 Hours", rotatefile.RotateTime(timex.OneHourSec*2).String())
	assert.Eq(t, "Every 15 Minutes", rotatefile.RotateTime(timex.OneMinSec*15).String())
	assert.Eq(t, "Every 5 Minutes", rotatefile.RotateTime(timex.OneMinSec*5).String())
	assert.Eq(t, "Every 3 Seconds", rotatefile.RotateTime(3).String())
}

func TestRotateTime_FirstCheckTime_Round(t *testing.T) {
	// log rotate interval minutes
	logMin := 5

	// now := timex.Now()
	// nowMin := now.Minute()
	nowMin := 37
	// dur := time.Duration(now.Minute() + min)
	dur := time.Duration(nowMin + logMin)
	assert.Eq(t, time.Duration(40), dur.Round(time.Duration(logMin)))

	nowMin = 40
	dur = time.Duration(nowMin + logMin)
	assert.Eq(t, time.Duration(45), dur.Round(time.Duration(logMin)))

	nowMin = 41
	dur = time.Duration(nowMin + logMin)
	assert.Eq(t, time.Duration(45), dur.Round(time.Duration(logMin)))
}
