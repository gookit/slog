package slog_test

import (
	"strings"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestFormattable_Format(t *testing.T) {
	r := newLogRecord("TEST_LOG_MESSAGE format")
	f := &slog.Formattable{}
	assert.Equal(t, "slog: TEST_LOG_MESSAGE format", r.GoString())

	bts, err := f.Format(r)
	assert.NoError(t, err)

	str := string(bts)
	assert.Contains(t, str, "TEST_LOG_MESSAGE format")

	fn := slog.FormatterFunc(func(r *slog.Record) ([]byte, error) {
		return []byte(r.Message), nil
	})

	bts, err = fn.Format(r)
	assert.NoError(t, err)

	str = string(bts)
	assert.Contains(t, str, "TEST_LOG_MESSAGE format")
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

func TestNewTextFormatter(t *testing.T) {
	f := slog.NewTextFormatter()

	dump.Println(f.Fields())
	assert.Contains(t, f.Fields(), "datetime")
	assert.Len(t, f.Fields(), strings.Count(slog.DefaultTemplate, "{{"))

	f.SetTemplate(slog.NamedTemplate)
	dump.Println(f.Fields())
	assert.Contains(t, f.Fields(), "datetime")
	assert.Len(t, f.Fields(), strings.Count(slog.NamedTemplate, "{{"))
}

func TestTextFormatter_Format(t *testing.T) {
	r := newLogRecord("TEST_LOG_MESSAGE")
	f := slog.NewTextFormatter()

	bs, err := f.Format(r)
	logTxt := string(bs)
	dump.Println(f.Template(), logTxt)

	assert.NoError(t, err)
	assert.NotEmpty(t, logTxt)
	assert.NotContains(t, logTxt, "{{")
	assert.NotContains(t, logTxt, "}}")
}

func TestJSONFormatter(t *testing.T) {
	l := slog.New()

	f := slog.NewJSONFormatter()
	f.AddField(slog.FieldKeyTimestamp)

	h := handler.NewConsoleHandler(slog.AllLevels)
	h.SetFormatter(f)

	l.AddHandler(h)

	fields := slog.M{
		"field1": 123,
		"field2": "abc",
	}

	l.WithFields(fields).Info("info", "message")

	// PrettyPrint=true

	l = slog.New()
	h = handler.NewConsoleHandler(slog.AllLevels)
	f = slog.NewJSONFormatter(func(f *slog.JSONFormatter) {
		f.Aliases = slog.StringMap{
			"level": "levelName",
		}
		f.PrettyPrint = true
	})

	h.SetFormatter(f)

	l.AddHandler(h)
	l.WithFields(fields).
		SetData(slog.M{"key1": "val1"}).
		SetExtra(slog.M{"ext1": "val1"}).
		Info("info message and PrettyPrint is TRUE")
}
