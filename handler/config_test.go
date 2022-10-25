package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func TestNewConfig(t *testing.T) {
	c := handler.NewConfig(
		handler.WithCompress(true),
	).With(handler.WithBuffSize(129))

	assert.True(t, c.Compress)
	assert.Eq(t, 129, c.BuffSize)
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

	h2 := handler.NewBuilder().
		WithOutput(new(bytes.Buffer)).
		WithUseJSON(true).
		Build()
	assert.NotNil(t, h2)

	assert.Panics(t, func() {
		handler.NewBuilder().Build()
	})
}
