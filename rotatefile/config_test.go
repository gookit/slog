package rotatefile_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/goutil/x/fmtutil"
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
	dump.P(cfg)

	assert.Eq(t, rotatefile.DefaultBackNum, cfg.BackupNum)
	assert.Eq(t, rotatefile.DefaultBackTime, cfg.BackupTime)
	assert.Eq(t, rotatefile.EveryHour, cfg.RotateTime)
	assert.Eq(t, rotatefile.DefaultMaxSize, cfg.MaxSize)
	assert.Eq(t, rotatefile.ModeRename, cfg.RotateMode)

	cfg = rotatefile.EmptyConfigWith(func(c *rotatefile.Config) {
		c.Compress = true
	})
	assert.True(t, cfg.Compress)
	assert.Eq(t, uint(0), cfg.BackupNum)
	assert.Eq(t, uint(0), cfg.BackupTime)

	cfg = &rotatefile.Config{}
	assert.Eq(t, rotatefile.ModeRename, cfg.RotateMode)

	err := jsonutil.DecodeString(`{
	"debug_mode": true,
	"rotate_mode": "create",
	"rotate_time": "1day"
}`, cfg)
	dump.P(cfg)
	assert.NoErr(t, err)
	assert.Eq(t, rotatefile.ModeCreate, cfg.RotateMode)
	assert.Eq(t, "Every 1 Day", cfg.RotateTime.String())
}

func TestRotateMode_cases(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Eq(t, "rename", rotatefile.ModeRename.String())
		assert.Eq(t, "create", rotatefile.ModeCreate.String())
		assert.Eq(t, "unknown", rotatefile.RotateMode(9).String())
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		rm := rotatefile.RotateMode(0)

		// UnmarshalJSON
		err := rm.UnmarshalJSON([]byte(`"create"`))
		assert.NoErr(t, err)
		assert.Eq(t, rotatefile.ModeCreate, rm)

		rm = rotatefile.RotateMode(0)
		// use int
		err = rm.UnmarshalJSON([]byte(`"1"`))
		assert.NoErr(t, err)
		assert.Eq(t, rotatefile.ModeCreate, rm)

		// error case
		assert.Err(t, rm.UnmarshalJSON([]byte(`create`)))
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		bs, err := rotatefile.ModeRename.MarshalJSON()
		assert.NoErr(t, err)
		assert.Eq(t, `"rename"`, string(bs))
		bs, err = rotatefile.ModeCreate.MarshalJSON()
		assert.NoErr(t, err)
		assert.Eq(t, `"create"`, string(bs))

		bs, err = rotatefile.RotateMode(35).MarshalJSON()
		assert.NoErr(t, err)
		assert.Eq(t, `"unknown"`, string(bs))
	})
}

func TestRotateTime_encode(t *testing.T) {
	rt := rotatefile.RotateTime(0)

	// UnmarshalJSON
	err := rt.UnmarshalJSON([]byte(`"1h"`))
	assert.NoErr(t, err)
	assert.Eq(t, "Every 1 Hours", rt.String())
	err = rt.UnmarshalJSON([]byte(`"3600"`))
	assert.NoErr(t, err)
	assert.Eq(t, "Every 1 Hours", rt.String())

	// error case
	assert.Err(t, rt.UnmarshalJSON([]byte(`a`)))

	// MarshalJSON
	bs, err := rt.MarshalJSON()
	assert.NoErr(t, err)
	assert.Eq(t, `"3600s"`, string(bs))
}

func TestRotateTime_TimeFormat(t *testing.T) {
	now := timex.Now()

	rt := rotatefile.EveryDay
	assert.Eq(t, "20060102", rt.TimeFormat())
	ft := rt.FirstCheckTime(now.T())
	assert.True(t, now.DayEnd().Equal(ft))

	rt = rotatefile.EveryHour
	assert.Eq(t, "20060102_1500", rt.TimeFormat())

	rt = rotatefile.Every15Min
	assert.Eq(t, "20060102_1504", rt.TimeFormat())
	ft = rt.FirstCheckTime(now.T())
	assert.Gt(t, ft.Unix(), 0)

	rt = rotatefile.EverySecond
	assert.Eq(t, "20060102_150405", rt.TimeFormat())
	ft = rt.FirstCheckTime(now.T())
	assert.Eq(t, now.Unix()+rt.Interval(), ft.Unix())
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
	assert.Eq(t, "Every 2 Day", rotatefile.RotateTime(timex.OneDaySec*2).String())
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
