package slog_test

import (
	"bytes"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestLoggerBasic(t *testing.T) {
	l := slog.New()
	l.SetName("testName")

	assert.Equal(t, "testName", l.Name())

	l = slog.NewWithName("testName")

	assert.Equal(t, "testName", l.Name())
}

func TestLogger_AddHandlers(t *testing.T) {

}

type bufferHandler struct {
	handler.BaseHandler
}

func (h *bufferHandler) Handle(*slog.Record) error {
	panic("implement me")
}

func TestLogger_ReportCaller(t *testing.T) {
	l := slog.NewWithConfig(func(logger *slog.Logger) {
		logger.ReportCaller = true
	})

	var buf bytes.Buffer
	h := handler.NewStreamHandler(&buf, slog.AllLevels)
	h.SetFormatter(slog.NewJSONFormatter(func(f *slog.JSONFormatter) {
		f.Fields = append(f.Fields, slog.FieldKeyCaller)
	}))

	l.AddHandler(h)
	l.Info("message")

	str := buf.String()
	assert.Contains(t, str, `"caller":"logger_test.go`)
}
