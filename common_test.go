package slog_test

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

var (
	testData1 = slog.M{"key0": "val0", "age": 23}
	testData2 = slog.M{"key0": "val0", "age": 23, "sub": slog.M{
		"subKey0": 345,
	}}
)

func TestDefine_basic(t *testing.T) {
	assert.NotEmpty(t, slog.NoTimeFields)
	assert.NotEmpty(t, slog.FieldKeyDate)
	assert.NotEmpty(t, slog.FieldKeyTime)
	assert.NotEmpty(t, slog.FieldKeyCaller)
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
	assert.Equal(t, "unknown", slog.Level(330).LowerName())
}

func TestLevelByName(t *testing.T) {
	assert.Equal(t, slog.InfoLevel, slog.LevelByName("info"))
	assert.Equal(t, slog.InfoLevel, slog.LevelByName("invalid"))
}

func TestLevelName(t *testing.T) {
	for level, wantName := range slog.LevelNames {
		realName := slog.LevelName(level)
		assert.Equal(t, wantName, realName)
	}

	assert.Equal(t, "UNKNOWN", slog.LevelName(20))
}

func TestLevel_ShouldHandling(t *testing.T) {
	assert.True(t, slog.InfoLevel.ShouldHandling(slog.ErrorLevel))
	assert.False(t, slog.InfoLevel.ShouldHandling(slog.TraceLevel))

	assert.True(t, slog.DebugLevel.ShouldHandling(slog.InfoLevel))
	assert.False(t, slog.DebugLevel.ShouldHandling(slog.TraceLevel))
}

func TestLevels_Contains(t *testing.T) {
	assert.True(t, slog.DangerLevels.Contains(slog.ErrorLevel))
	assert.False(t, slog.DangerLevels.Contains(slog.InfoLevel))
	assert.True(t, slog.NormalLevels.Contains(slog.InfoLevel))
	assert.False(t, slog.NormalLevels.Contains(slog.PanicLevel))
}

func TestNewLvFormatter(t *testing.T) {
	lf := slog.NewLvFormatter(slog.InfoLevel)

	assert.True(t, lf.IsHandling(slog.ErrorLevel))
	assert.True(t, lf.IsHandling(slog.InfoLevel))
	assert.False(t, lf.IsHandling(slog.DebugLevel))
}
