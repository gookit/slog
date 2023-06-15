package rotatefile_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/rotatefile"
)

func TestNewFilesClear(t *testing.T) {
	fc := rotatefile.NewFilesClear(nil)
	fc.WithConfigFn(func(c *rotatefile.CConfig) {
		c.AddFileDir("./testdata")
	})

	err := fc.Clean()
	assert.NoErr(t, err)
}
