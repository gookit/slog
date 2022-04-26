package slog_test

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func Example_quickStart() {
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Infof("info log %s", "message")
	slog.Debugf("debug %s", "message")
}

func Example_configSlog() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	slog.Trace("this is a simple log message")
	slog.Debug("this is a simple log message")
	slog.Info("this is a simple log message")
	slog.Notice("this is a simple log message")
	slog.Warn("this is a simple log message")
	slog.Error("this is a simple log message")
	slog.Fatal("this is a simple log message")
}

func Example_useJSONFormat() {
	// use JSON formatter
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

func ExampleNew() {
	mylog := slog.New()
	levels := slog.AllLevels

	mylog.AddHandler(handler.MustFileHandler("app.log", handler.WithLogLevels(levels)))

	mylog.Info("info log message")
	mylog.Warn("warning log message")
	mylog.Infof("info log %s", "message")
}
