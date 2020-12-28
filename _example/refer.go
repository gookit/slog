package main

import (
	"flag"
	"log"
	"time"

	"github.com/golang/glog"
	"github.com/gookit/slog"
	"github.com/sirupsen/logrus"

	"github.com/syyongx/llog"

	"go.uber.org/zap"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	// for glog
	flag.Parse()

	// -- log
	log.Println("raw log message")

	// -- glog
	glog.Infof("glog %s", "message message")

	// -- llog
	llog.NewLogger("llog test").Info("llog message message")

	// -- slog
	slog.Debug("slog message message")
	slog.WithFields(slog.M{
		"omg":    true,
		"number": 122,
	}).Infof("slog %s", "message message")

	// -- logrus
	logrus.Debug("logrus message message")
	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	// -- zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlog.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("zerolog message")
	zlog.Print("zerolog hello")

	// slog.Infof("log %s", "message")
	url := "/path/to/some"

	// -- zap
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
