package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/x/fmtutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func TestNewConfig(t *testing.T) {
	c := handler.NewConfig(
		handler.WithCompress(true),
		handler.WithLevelMode(handler.LevelModeValue),
		handler.WithBackupNum(20),
		handler.WithBackupTime(1800),
		handler.WithRotateMode(rotatefile.ModeCreate),
		func(c *handler.Config) {
			c.BackupTime = 23
			c.RenameFunc = func(fpath string, num uint) string {
				return fpath + ".bak"
			}
		},
	).
		With(handler.WithBuffSize(129)).
		WithConfigFn(handler.WithLogLevel(slog.ErrorLevel))

	assert.True(t, c.Compress)
	assert.Eq(t, 129, c.BuffSize)
	assert.Eq(t, handler.LevelModeValue, c.LevelMode)
	assert.Eq(t, slog.ErrorLevel, c.Level)
	assert.Eq(t, rotatefile.ModeCreate, c.RotateMode)

	c.With(handler.WithLevelModeString("max"))
	assert.Eq(t, slog.LevelModeMax, c.LevelMode)

	c.WithConfigFn(handler.WithLevelNames([]string{"info", "debug"}))
	assert.Eq(t, []slog.Level{slog.InfoLevel, slog.DebugLevel}, c.Levels)
}

func TestConfig_fromJSON(t *testing.T) {
	c := &handler.Config{}
	assert.Eq(t, slog.LevelModeList, c.LevelMode)
	assert.Eq(t, rotatefile.ModeRename, c.RotateMode)

	assert.NoErr(t, c.FromJSON([]byte(`{
		"logfile": "testdata/config_test.log",
		"level": "debug",
		"level_mode": "max",
		"levels": ["info", "debug"],
		"buff_mode": "line",
		"buff_size": 128,
		"backup_num": 3,
		"backup_time": 3600,
		"rotate_mode": "create",
		"rotate_time": "1day"
	}`)))
	c.With(handler.WithDebugMode)
	dump.P(c)

	assert.Eq(t, slog.LevelModeMax, c.LevelMode)
	assert.Eq(t, rotatefile.ModeCreate, c.RotateMode)
	assert.Eq(t, "Every 1 Day", c.RotateTime.String())
}

func TestWithLevelNamesString(t *testing.T) {
	c := handler.NewConfig(handler.WithLevelNamesString("info,error"))
	assert.Eq(t, []slog.Level{slog.InfoLevel, slog.ErrorLevel}, c.Levels)
}

func TestWithMaxLevelName(t *testing.T) {
	c := handler.NewConfig(handler.WithMaxLevelName("error"))
	assert.Eq(t, slog.ErrorLevel, c.Level)
	assert.Eq(t, handler.LevelModeValue, c.LevelMode)

	c1 := handler.NewConfig(handler.WithLevelName("warn"))
	assert.Eq(t, slog.WarnLevel, c1.Level)
	assert.Eq(t, handler.LevelModeValue, c1.LevelMode)
}

func TestWithRotateMode(t *testing.T) {
	c := handler.Config{}

	c.With(handler.WithRotateModeString("rename"))
	assert.Eq(t, rotatefile.ModeRename, c.RotateMode)

	assert.PanicsErrMsg(t, func() {
		c.With(handler.WithRotateModeString("invalid"))
	}, "rotatefile: invalid rotate mode: invalid")

}

func TestWithRotateTimeString(t *testing.T) {
	tests := []struct {
		input    string
		expected rotatefile.RotateTime
		panics   bool
	}{
		{"1hours", rotatefile.RotateTime(3600), false},
		{"24h", rotatefile.RotateTime(86400), false},
		{"1day", rotatefile.RotateTime(86400), false},
		{"7d", rotatefile.RotateTime(604800), false},
		{"1m", rotatefile.RotateTime(60), false},
		{"30s", rotatefile.RotateTime(30), false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			c := &handler.Config{}
			if tt.panics {
				assert.Panics(t, func() {
					handler.WithRotateTimeString(tt.input)(c)
				})
			} else {
				assert.NotPanics(t, func() {
					handler.WithRotateTimeString(tt.input)(c)
				})
				assert.Eq(t, tt.expected, c.RotateTime)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	testFile := "testdata/builder.log"
	assert.NoErr(t, fsutil.DeleteIfFileExist(testFile))

	b := handler.NewBuilder().
		WithLogfile(testFile).
		WithLogLevels(slog.AllLevels).
		WithBuffSize(128).
		WithBuffMode(handler.BuffModeBite).
		WithMaxSize(fmtutil.OneMByte * 3).
		WithRotateTime(rotatefile.Every30Min).
		WithCompress(true).
		With(func(c *handler.Config) {
			c.BackupNum = 3
		})

	assert.Eq(t, uint(3), b.BackupNum)
	assert.Eq(t, handler.BuffModeBite, b.BuffMode)
	assert.Eq(t, rotatefile.Every30Min, b.RotateTime)

	h := b.Build()
	assert.NotNil(t, h)
	assert.NoErr(t, h.Close())

	b1 := handler.NewBuilder().
		WithOutput(new(bytes.Buffer)).
		WithUseJSON(true).
		WithLogLevel(slog.ErrorLevel).
		WithLevelMode(handler.LevelModeValue)
	assert.Eq(t, handler.LevelModeValue, b1.LevelMode)
	assert.Eq(t, slog.ErrorLevel, b1.Level)

	h2 := b1.Build()
	assert.NotNil(t, h2)

	assert.Panics(t, func() {
		handler.NewBuilder().Build()
	})
}

type simpleWriter struct {
	errOnWrite bool
}

func (w *simpleWriter) Write(p []byte) (n int, err error) {
	if w.errOnWrite {
		return 0, errorx.Raw("write error")
	}
	return len(p), nil
}

type closeWriter struct {
	errOnWrite bool
	errOnClose bool
}

func (w *closeWriter) Close() error {
	if w.errOnClose {
		return errorx.Raw("close error")
	}
	return nil
}

func (w *closeWriter) Write(p []byte) (n int, err error) {
	if w.errOnWrite {
		return 0, errorx.Raw("write error")
	}
	return len(p), nil
}

type flushCloseWriter struct {
	closeWriter
	errOnFlush bool
}

// Flush implement stdio.Flusher
func (w *flushCloseWriter) Flush() error {
	if w.errOnFlush {
		return errorx.Raw("flush error")
	}
	return nil
}

type syncCloseWriter struct {
	closeWriter
	errOnSync bool
}

// Sync implement stdio.Syncer
func (w *syncCloseWriter) Sync() error {
	if w.errOnSync {
		return errorx.Raw("sync error")
	}
	return nil
}

func TestNewBuilder_buildFromWriter(t *testing.T) {
	t.Run("FlushCloseWriter", func(t *testing.T) {
		out := &flushCloseWriter{}
		out.errOnFlush = true
		h := handler.NewBuilder().
			WithOutput(out).
			WithConfigFn(func(c *handler.Config) {
				c.RenameFunc = func(fpath string, num uint) string {
					return fpath + ".bak"
				}
			}).
			Build()
		assert.Err(t, h.Flush())

		// wrap buffer
		h = handler.NewBuilder().
			WithOutput(out).
			WithBuffSize(128).
			Build()
		assert.NoErr(t, h.Close())
		assert.NoErr(t, h.Flush())
	})

	t.Run("CloseWriter", func(t *testing.T) {
		h := handler.NewBuilder().
			WithOutput(&closeWriter{errOnClose: true}).
			WithBuffSize(128).
			Build()
		assert.NotNil(t, h)
		assert.Err(t, h.Close())
	})

	t.Run("SimpleWriter", func(t *testing.T) {
		h := handler.NewBuilder().
			WithOutput(&simpleWriter{errOnWrite: true}).
			WithBuffSize(128).
			Build()
		assert.NotNil(t, h)
		assert.NoErr(t, h.Close())
	})
}
