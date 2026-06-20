package slog_test

import (
	"io"
	"sync/atomic"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

type countStringer struct{ n *int32 }

func (c countStringer) String() string { atomic.AddInt32(c.n, 1); return "X" }

// regression: a disabled level must not build/format the message at all.
func TestLogger_DisabledLevel_SkipsFormatting(t *testing.T) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.ErrorLevel}),
	)

	var cnt int32
	for i := 0; i < 100; i++ {
		logger.Info("msg", countStringer{&cnt}) // Info is filtered out
	}
	assert.Eq(t, int32(0), atomic.LoadInt32(&cnt))

	// the enabled level still formats
	logger.Error("msg", countStringer{&cnt})
	assert.Eq(t, int32(1), atomic.LoadInt32(&cnt))
}

// Panic/Fatal must still trigger their side effects even with no matching handler.
func TestLogger_PanicFatal_AlwaysHandled_NoMatchingHandler(t *testing.T) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.InfoLevel}),
	)

	panicked := false
	logger.PanicFunc = func(v any) { panicked = true }
	logger.Panic("boom")
	assert.True(t, panicked)
}
