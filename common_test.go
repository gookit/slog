package slog_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/gsr"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
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
		"k6": map[string]any{
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

func newLogRecord(msg string) *slog.Record {
	r := &slog.Record{
		Channel: slog.DefaultChannelName,
		Level:   slog.InfoLevel,
		Message: msg,
		Time:    slog.DefaultClockFn.Now(),
		Data: map[string]any{
			"data_key0": "value",
			"username":  "inhere",
		},
		Extra: map[string]any{
			"source":     "linux",
			"extra_key0": "hello",
		},
		// Caller: goinfo.GetCallerInfo(),
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
	slog.FormatterWrapper
	byteutil.Buffer
	errOnHandle bool
	errOnClose  bool
	errOnFlush  bool
	// hooks
	callOnFlush func()
	// tip: 设置为true，默认会让 error,fatal 等信息提前被reset丢弃掉.
	// see Logger.writeRecord()
	resetOnFlush bool
}

// built in test, will collect logs to buffer
func newTestHandler() *testHandler {
	return &testHandler{resetOnFlush: true}
}

func (h *testHandler) IsHandling(_ slog.Level) bool {
	return true
}

func (h *testHandler) Close() error {
	if h.errOnClose {
		return errorx.Raw("close error")
	}

	h.Reset()
	return nil
}

func (h *testHandler) Flush() error {
	if h.errOnFlush {
		return errorx.Raw("flush error")
	}
	if h.callOnFlush != nil {
		h.callOnFlush()
	}

	if h.resetOnFlush {
		h.Reset()
	}
	return nil
}

func (h *testHandler) Handle(r *slog.Record) error {
	if h.errOnHandle {
		return errorx.Raw("handle error")
	}

	bs, err := h.Format(r)
	if err != nil {
		return err
	}
	h.Write(bs)
	return nil
}

type testFormatter struct {
	errOnFormat bool
}

func newTestFormatter(errOnFormat ...bool) *testFormatter {
	return &testFormatter{
		errOnFormat: len(errOnFormat) > 0 && errOnFormat[0],
	}
}

func (f testFormatter) Format(r *slog.Record) ([]byte, error) {
	if f.errOnFormat {
		return nil, errorx.Raw("format error")
	}
	return []byte(r.Message), nil
}

func newLogger() *slog.Logger {
	return slog.NewWithConfig(func(l *slog.Logger) {
		l.ReportCaller = true
		l.DoNothingOnPanicFatal()
	})
}

// newTestLogger create a logger for test, will write logs to buffer
func newTestLogger() (*closedBuffer, *slog.Logger) {
	l := slog.NewWithConfig(func(l *slog.Logger) {
		l.DoNothingOnPanicFatal()
		l.CallerFlag = slog.CallerFlagFull
	})
	w := newBuffer()
	h := handler.NewIOWriter(w, slog.AllLevels)
	// fmt.Print("Template:", h.TextFormatter().Template())
	l.SetHandlers([]slog.Handler{h})
	return w, l
}

func printAllLevelLogs(l gsr.Logger, args ...any) {
	l.Debug(args...)
	l.Info(args...)
	l.Warn(args...)
	l.Error(args...)
	l.Print(args...)
	l.Println(args...)
	l.Fatal(args...)
	l.Fatalln(args...)
	l.Panic(args...)
	l.Panicln(args...)

	sl, ok := l.(*slog.Logger)
	if ok {
		sl.Trace(args...)
		sl.Notice(args...)
		sl.ErrorT(errorx.Raw("a error object"))
		sl.ErrorT(errorx.New("error with stack info"))
	}
}

func printfAllLevelLogs(l gsr.Logger, tpl string, args ...any) {
	l.Printf(tpl, args...)
	l.Debugf(tpl, args...)
	l.Infof(tpl, args...)
	l.Warnf(tpl, args...)
	l.Errorf(tpl, args...)
	l.Panicf(tpl, args...)
	l.Fatalf(tpl, args...)

	if sl, ok := l.(*slog.Logger); ok {
		sl.Noticef(tpl, args...)
		sl.Tracef(tpl, args...)
	}
}
