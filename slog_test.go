package slog_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/formatter"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func TestInfof(t *testing.T) {
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(formatter.NewJSONFormatter(slog.StringMap{
		"level": "levelName",
		"message": "msg",
		"data": "params",
	}))

	slog.AddHandler(h2)
	slog.AddProcessor(slog.AddHostname())

	slog.Infof("info %s", "message")
}

func TestName2Level(t *testing.T) {
	for wantLevel, name := range slog.LevelNames {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	level, err := slog.Name2Level("unknown")
	assert.Error(t, err)
	assert.Equal(t, slog.Level(0), level)
}
