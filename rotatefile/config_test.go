package rotatefile_test

import (
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/slog/rotatefile"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	size := fmtutil.DataSize(1024 * 1024 * 10)
	dump.P(size)

	c := rotatefile.NewDefaultConfig()
	assert.Equal(t, rotatefile.DefaultMaxSize, c.MaxSize)
}
