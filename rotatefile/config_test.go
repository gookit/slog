package rotatefile_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog/rotatefile"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	size := fmtutil.DataSize(1024 * 1024 * 10)
	dump.P(size)

	c := rotatefile.NewDefaultConfig()
	assert.Equal(t, rotatefile.DefaultMaxSize, c.MaxSize)
}

func TestNewConfig(t *testing.T) {
	cfg := rotatefile.NewConfig("testdata/test.log")

	assert.Equal(t, rotatefile.DefaultBackNum, cfg.BackupNum)
	assert.Equal(t, rotatefile.DefaultBackTime, cfg.BackupTime)
	assert.Equal(t, rotatefile.EveryHour, cfg.RotateTime)
	assert.Equal(t, rotatefile.DefaultMaxSize, cfg.MaxSize)

	dump.P(cfg)
}

func TestRotateTime_String(t *testing.T) {
	assert.Equal(t, "Every 1 Day", rotatefile.EveryDay.String())
	assert.Equal(t, "Every 1 Hours", rotatefile.EveryHour.String())
	assert.Equal(t, "Every 1 Minutes", rotatefile.EveryMinute.String())
	assert.Equal(t, "Every 1 Seconds", rotatefile.EverySecond.String())

	assert.Equal(t, "Every 2 Hours", rotatefile.RotateTime(timex.OneHourSec*2).String())
	assert.Equal(t, "Every 15 Minutes", rotatefile.RotateTime(timex.OneMinSec*15).String())
	assert.Equal(t, "Every 5 Minutes", rotatefile.RotateTime(timex.OneMinSec*5).String())
	assert.Equal(t, "Every 3 Seconds", rotatefile.RotateTime(3).String())
}

func TestRotateTime_FirstCheckTime_Round(t *testing.T) {
	// log rotate interval minutes
	logMin := 5

	// now := timex.Now()
	// nowMin := now.Minute()
	nowMin := 37
	// dur := time.Duration(now.Minute() + min)
	dur := time.Duration(nowMin + logMin)
	assert.Equal(t, time.Duration(40), dur.Round(time.Duration(logMin)))

	nowMin = 40
	dur = time.Duration(nowMin + logMin)
	assert.Equal(t, time.Duration(45), dur.Round(time.Duration(logMin)))

	nowMin = 41
	dur = time.Duration(nowMin + logMin)
	assert.Equal(t, time.Duration(45), dur.Round(time.Duration(logMin)))
}
