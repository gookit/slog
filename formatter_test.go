package slog_test

import (
	"strings"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func TestFormattableTrait_Formatter(t *testing.T) {
	ft := &slog.FormattableTrait{}
	tf := slog.AsTextFormatter(ft.Formatter())
	assert.NotNil(t, tf)
	assert.Panics(t, func() {
		slog.AsJSONFormatter(ft.Formatter())
	})

	ft.SetFormatter(slog.NewJSONFormatter())
	jf := slog.AsJSONFormatter(ft.Formatter())
	assert.NotNil(t, jf)
	assert.Panics(t, func() {
		slog.AsTextFormatter(ft.Formatter())
	})
}

func TestFormattable_Format(t *testing.T) {
	r := newLogRecord("TEST_LOG_MESSAGE format")
	f := &slog.FormattableTrait{}
	assert.Eq(t, "slog: TEST_LOG_MESSAGE format", r.GoString())

	bts, err := f.Format(r)
	assert.NoErr(t, err)

	str := string(bts)
	assert.Contains(t, str, "TEST_LOG_MESSAGE format")

	fn := slog.FormatterFunc(func(r *slog.Record) ([]byte, error) {
		return []byte(r.Message), nil
	})

	bts, err = fn.Format(r)
	assert.NoErr(t, err)

	str = string(bts)
	assert.Contains(t, str, "TEST_LOG_MESSAGE format")
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

	assert.NoErr(t, err)
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
