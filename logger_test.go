package slog_test

import (
	"testing"

	"github.com/gookit/slog"
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
