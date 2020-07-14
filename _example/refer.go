package main

import (
	"flag"
	"time"

	"github.com/golang/glog"

	// "github.com/gookit/slog"
	"github.com/sirupsen/logrus"
	"github.com/syyongx/llog"

	"go.uber.org/zap"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	// for glog
	flag.Parse()

	logrus.Debug("logrus message")
	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	glog.Infof("glog %s", "message")

	llog.NewLogger("llog test")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zlog.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("zerolog message")
	zlog.Print("zerolog hello")

	// slog.Infof("log %s", "message")
	url := "/path/to/some"

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("zap log. Failed to fetch URL: %s", url)
}
