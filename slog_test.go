package slog_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

var doNothing = func(code int) {
	// do nothing
}

func TestStd(t *testing.T) {
	defer slog.Reset()
	assert.Equal(t, "stdLogger", slog.Std().Name())

	_, ok := slog.GetFormatter().(*slog.TextFormatter)
	assert.True(t, ok)

	slog.SetFormatter(slog.NewJSONFormatter())
	_, ok = slog.GetFormatter().(*slog.JSONFormatter)
	assert.True(t, ok)

	buf := new(bytes.Buffer)
	slog.Std().ExitFunc = func(code int) {
		buf.WriteString("Exited,")
		buf.WriteString(strconv.Itoa(code))
	}
	slog.Exit(34)
	assert.Equal(t, "Exited,34", buf.String())
}

func TestTextFormatNoColor(t *testing.T) {
	defer slog.Reset()
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = false
	})

	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Infof("info log %s", "message")
	slog.Debugf("debug %s", "message")
	slog.Reset()
}

func TestTextFormatWithColor(t *testing.T) {
	defer slog.Reset()

	slog.Configure(func(l *slog.SugaredLogger) {
		l.Level = slog.PanicLevel
		l.ExitFunc = doNothing
	})

	printLogs()

	slog.GetFormatter().(*slog.TextFormatter).Template = slog.NamedTemplate

	printLogs()
}

func printLogs() {
	slog.Std().Trace("this is a simple log message")
	slog.Trace("this is a simple log message")
	slog.Debug("this is a simple log message")
	slog.Info("this is a simple log message")
	slog.Notice("this is a simple log message")
	slog.Warn("this is a simple log message")
	slog.Error("this is a simple log message")
	slog.Fatal("this is a simple log message")
}

func TestUseJSONFormat(t *testing.T) {
	defer slog.Reset()
	slog.SetFormatter(slog.NewJSONFormatter())

	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.WithData(slog.M{
		"key0": 134,
		"key1": "abc",
	}).Infof("info log %s", "message")

	r := slog.WithFields(slog.M{
		"category": "service",
		"IP":       "127.0.0.1",
	})
	r.Infof("info %s", "message")
	r.Debugf("debug %s", "message")
}

func TestAddHandler(t *testing.T) {
	defer slog.Reset()
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(slog.NewJSONFormatter().Configure(func(f *slog.JSONFormatter) {
		f.Aliases = slog.StringMap{
			"level":   "levelName",
			"message": "msg",
			"data":    "params",
		}
	}))

	slog.AddHandlers(h2)
	slog.Infof("info %s", "message")
}

func TestAddProcessor(t *testing.T) {
	defer slog.Reset()

	buf := new(bytes.Buffer)
	slog.Configure(func(logger *slog.SugaredLogger) {
		logger.Output = buf
		logger.Formatter = slog.NewJSONFormatter()
	})

	slog.AddProcessor(slog.AddHostname())
	slog.Trace("Trace message")
	slog.Tracef("Tracef %s", "message")

	str := buf.String()
	buf.Reset()
	fmt.Println(str)
	assert.Contains(t, str, `"hostname":`)
	assert.Contains(t, str, "Trace message")
	assert.Contains(t, str, "Tracef message")

	slog.AddProcessors(slog.ProcessorFunc(func(r *slog.Record) {
		r.AddField("newField", "newValue")
	}))
	slog.Debug("Debug message")
	slog.Debugf("Debugf %s", "message")
	str = buf.String()
	buf.Reset()

	assert.Contains(t, str, `"newField":"newValue"`)
	assert.Contains(t, str, "Debug message")
	assert.Contains(t, str, "Debugf message")
}

func TestLevelName(t *testing.T) {
	for level, wantName := range slog.LevelNames {
		realName := slog.LevelName(level)
		assert.Equal(t, wantName, realName)
	}

	assert.Equal(t, "UNKNOWN", slog.LevelName(20))
}

func TestName2Level(t *testing.T) {
	for wantLevel, name := range slog.LevelNames {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	// special names
	tests := map[slog.Level]string{
		slog.WarnLevel:  "warn",
		slog.ErrorLevel: "err",
		slog.InfoLevel:  "",
	}
	for wantLevel, name := range tests {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	level, err := slog.Name2Level("unknown")
	assert.Error(t, err)
	assert.Equal(t, slog.Level(0), level)
}
