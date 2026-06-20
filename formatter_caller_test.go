package slog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// regression: with ReportCaller=false and the default template (which has
// {{caller}}), the text formatter must render an empty caller, not the literal
// string "caller".
func TestTextFormatter_NilCaller_NotLiteral(t *testing.T) {
	buf := &bytes.Buffer{}
	l := slog.NewWithHandlers(handler.NewIOWriter(buf, slog.AllLevels))
	l.ReportCaller = false // default template still contains {{caller}}

	l.Info("hello")
	out := buf.String()
	assert.StrContains(t, out, "hello")
	assert.False(t, strings.Contains(out, "caller"), "should not print literal 'caller'")
}
