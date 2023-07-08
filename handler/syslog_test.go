//go:build !windows && !plan9

package handler_test

import (
	"log/syslog"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/handler"
)

func TestNewSysLogHandler(t *testing.T) {
	h, err := handler.NewSysLogHandler(syslog.LOG_INFO, "slog")
	assert.NoErr(t, err)

	err = h.Handle(newLogRecord("test syslog handler"))
	assert.NoErr(t, err)

	assert.NoErr(t, h.Flush())
	assert.NoErr(t, h.Close())
}
