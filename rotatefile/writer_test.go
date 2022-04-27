package rotatefile_test

import (
	"testing"

	"github.com/gookit/slog/rotatefile"
	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	wr, err := rotatefile.NewConfig("testdata/test.log").Create()
	if err != nil {
		return
	}

	c := wr.Config()
	assert.Equal(t, c.MaxSize, rotatefile.DefaultMaxSize)
}
