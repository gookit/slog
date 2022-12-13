package slog_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
)

var (
	testData1 = slog.M{"key0": "val0", "age": 23}
	// testData2 = slog.M{"key0": "val0", "age": 23, "sub": slog.M{
	// 	"subKey0": 345,
	// }}
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
	assert.Eq(t, "INFO", slog.InfoLevel.Name())
	assert.Eq(t, "INFO", slog.InfoLevel.String())
	assert.Eq(t, "info", slog.InfoLevel.LowerName())
	assert.Eq(t, "unknown", slog.Level(330).LowerName())
}

func TestLevelByName(t *testing.T) {
	assert.Eq(t, slog.InfoLevel, slog.LevelByName("info"))
	assert.Eq(t, slog.InfoLevel, slog.LevelByName("invalid"))
}

func TestLevelName(t *testing.T) {
	for level, wantName := range slog.LevelNames {
		realName := slog.LevelName(level)
		assert.Eq(t, wantName, realName)
	}

	assert.Eq(t, "UNKNOWN", slog.LevelName(20))
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

func newLogRecord(msg string) *slog.Record {
	r := &slog.Record{
		Channel: slog.DefaultChannelName,
		Level:   slog.InfoLevel,
		Message: msg,
		Time:    slog.DefaultClockFn.Now(),
		Data: map[string]interface{}{
			"data_key0": "value",
			"username":  "inhere",
		},
		Extra: map[string]interface{}{
			"source":     "linux",
			"extra_key0": "hello",
		},
		// Caller: stdutil.GetCallerInfo(),
	}

	r.Init(true)
	return r
}

type closedBuffer struct {
	bytes.Buffer
}

func newBuffer() *closedBuffer {
	return &closedBuffer{}
}

func (w *closedBuffer) Close() error {
	return nil
}

func (w *closedBuffer) StringReset() string {
	s := w.Buffer.String()
	w.Reset()
	return s
}

type testHandler struct {
	errOnHandle bool
	errOnFlush  bool
	errOnClose  bool
}

func newTestHandler() *testHandler {
	return &testHandler{}
}

// func (h testHandler) Reset() {
// 	h.errOnHandle = false
// 	h.errOnFlush = false
// 	h.errOnClose = false
// }

func (h testHandler) IsHandling(_ slog.Level) bool {
	return true
}

func (h testHandler) Close() error {
	if h.errOnClose {
		return errorx.Raw("close error")
	}
	return nil
}

func (h testHandler) Flush() error {
	if h.errOnFlush {
		return errorx.Raw("flush error")
	}
	return nil
}

func (h testHandler) Handle(_ *slog.Record) error {
	if h.errOnHandle {
		return errorx.Raw("handle error")
	}
	return nil
}
