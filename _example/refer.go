package main

import (
	"github.com/golang/glog"

	// "github.com/gookit/slog"
	"github.com/sirupsen/logrus"
	"github.com/syyongx/llog"
)

func main() {
	logrus.Debug("message")
	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	glog.Infof("log %s", "message")

	llog.NewLogger("test")

	// slog.Infof("log %s", "message")
}
