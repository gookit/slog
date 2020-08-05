package slog_test

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

func TestPrependExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	assert.Len(t, slog.ExitHandlers(), 2)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Equal(t, "HANDLER2-HANDLER1-Exited", buf.String())
}

func TestRegisterExitHandler(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	buf := new(bytes.Buffer)
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER1-")
	})
	slog.RegisterExitHandler(func() {
		buf.WriteString("HANDLER2-")
	})
	// prepend
	slog.PrependExitHandler(func() {
		buf.WriteString("HANDLER3-")
	})
	assert.Len(t, slog.ExitHandlers(), 3)

	slog.SetExitFunc(func(code int) {
		buf.WriteString("Exited")
	})
	slog.Exit(23)
	assert.Equal(t, "HANDLER3-HANDLER1-HANDLER2-Exited", buf.String())
}

func TestExitHandlerWithError(t *testing.T) {
	defer slog.Reset()

	assert.Len(t, slog.ExitHandlers(), 0)

	slog.RegisterExitHandler(func() {
		panic("test error")
	})

	slog.SetExitFunc(func(code int) {})

	testutil.RewriteStderr()
	slog.Exit(23)
	str := testutil.RestoreStderr()
	assert.Equal(t, "Run exit handler error: test error\n", str)
}

func TestLogger_ExitHandlerWithError(t *testing.T) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.ExitFunc = doNothing
	})

	assert.Len(t, l.ExitHandlers(), 0)

	l.RegisterExitHandler(func() {
		panic("test error")
	})

	testutil.RewriteStderr()
	l.Exit(23)
	str := testutil.RestoreStderr()
	assert.Equal(t, "Run exit handler error: test error\n", str)
}
