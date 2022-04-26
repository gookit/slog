package slog_test

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

func TestDefine_basic(t *testing.T) {
	assert.NotEmpty(t, slog.NoTimeFields)
	assert.NotEmpty(t, slog.FieldKeyDate)
	assert.NotEmpty(t, slog.FieldKeyTime)
	assert.NotEmpty(t, slog.FieldKeyPkg)
	assert.NotEmpty(t, slog.FieldKeyError)
}

func TestM_String(t *testing.T) {
	m := slog.M{
		"k0": 12,
		"k1": "abc",
		"k2": true,
		"k3": 23.45,
		"k4": []int{12, 23},
		"k5": []string{"ab", "bc"},
		"k6": map[string]interface{}{
			"k6-1": 23,
			"k6-2": "def",
		},
	}

	fmt.Println(m)
	dump.P(m.String(), m)
	assert.NotEmpty(t, m.String())
}

func TestLevel_Name(t *testing.T) {
	assert.Equal(t, "INFO", slog.InfoLevel.Name())
	assert.Equal(t, "INFO", slog.InfoLevel.String())
	assert.Equal(t, "info", slog.InfoLevel.LowerName())
}

func TestLevel_ShouldHandling(t *testing.T) {
	assert.True(t, slog.InfoLevel.ShouldHandling(slog.ErrorLevel))
	assert.False(t, slog.InfoLevel.ShouldHandling(slog.TraceLevel))

	assert.True(t, slog.DebugLevel.ShouldHandling(slog.InfoLevel))
	assert.False(t, slog.DebugLevel.ShouldHandling(slog.TraceLevel))
}
