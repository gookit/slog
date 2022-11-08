package main

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/phuslu/log"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// In _example/ dir, run:
//
//	go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_loglibs_test.go
//
// code refer:
//
//	https://github.com/phuslu/log
var msg = "The quick brown fox jumps over the lazy dog"

func BenchmarkZapNegative(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(ioutil.Discard),
		zapcore.ErrorLevel,
	))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg, zap.String("rate", "15"), zap.Int("low", 16), zap.Float32("high", 123.2))
	}
}

func BenchmarkZeroLogNegative(b *testing.B) {
	logger := zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkPhusLogNegative(b *testing.B) {
	logger := log.Logger{Level: log.ErrorLevel, Writer: log.IOWriter{Writer: ioutil.Discard}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

// "github.com/sirupsen/logrus"
func BenchmarkLogrusNegative(b *testing.B) {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	logger.Level = logrus.ErrorLevel

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkGookitSlogNegative(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(ioutil.Discard, []slog.Level{slog.ErrorLevel}),
	)
	logger.ReportCaller = false

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkZapPositive(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(ioutil.Discard),
		zapcore.InfoLevel,
	))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg, zap.String("rate", "15"), zap.Int("low", 16), zap.Float32("high", 123.2))
	}
}

func BenchmarkZeroLogPositive(b *testing.B) {
	logger := zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkPhusLogPositive(b *testing.B) {
	logger := log.Logger{Writer: log.IOWriter{Writer: ioutil.Discard}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkLogrusPositive(b *testing.B) {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	logger.Level = logrus.TraceLevel

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkGookitSlogPositive(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(ioutil.Discard, slog.NormalLevels),
	)
	logger.ReportCaller = false

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}
