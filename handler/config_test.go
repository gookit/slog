package handler_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := handler.NewConfig(
		handler.WithCompress(true),
	).With(handler.WithBuffSize(129))

	assert.True(t, c.Compress)
	assert.Equal(t, 129, c.BuffSize)
}

func TestNewBuilder(t *testing.T) {
	testFile := "testdata/builder.log"
	assert.NoError(t, fsutil.DeleteIfFileExist(testFile))

	h := handler.NewBuilder().
		With(
			handler.WithLogfile(testFile),
			handler.WithBuffSize(128),
			handler.WithBuffMode(handler.BuffModeBite),
			handler.WithCompress(true),
		).
		Build()
	assert.NotNil(t, h)
	assert.NoError(t, h.Close())

	h2 := handler.NewBuilder().
		WithOutput(new(bytes.Buffer)).
		With(handler.WithUseJSON(true)).
		Build()
	assert.NotNil(t, h2)

	assert.Panics(t, func() {
		handler.NewBuilder().Build()
	})
}
