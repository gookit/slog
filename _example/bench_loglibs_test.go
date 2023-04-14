package main

import (
	"io"
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	phuslu "github.com/phuslu/log"
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
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg, zap.String("rate", "15"), zap.Int("low", 16), zap.Float32("high", 123.2))
	}
}

func BenchmarkZapSugarNegative(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		// zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)).Sugar()

	// logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	// return

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkZeroLogNegative(b *testing.B) {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger().Level(zerolog.InfoLevel)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkPhusLogNegative(b *testing.B) {
	logger := phuslu.Logger{Level: phuslu.InfoLevel, Writer: phuslu.IOWriter{Writer: io.Discard}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

// "github.com/sirupsen/logrus"
func BenchmarkLogrusNegative(b *testing.B) {
	logger := logrus.New()
	logger.Out = io.Discard
	logger.Level = logrus.InfoLevel

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkGookitSlogNegative(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.InfoLevel}),
		// handler.NewIOWriter(os.Stdout, []slog.Level{slog.InfoLevel}),
	)
	logger.ReportCaller = false

	// logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	// return

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkZapPositive(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg, zap.String("rate", "15"), zap.Int("low", 16), zap.Float32("high", 123.2))
	}
}

func BenchmarkZapSugarPositive(b *testing.B) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)).Sugar()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg, zap.String("rate", "15"), zap.Int("low", 16), zap.Float32("high", 123.2))
	}
}

func BenchmarkZeroLogPositive(b *testing.B) {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger().Level(zerolog.InfoLevel)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkPhusLogPositive(b *testing.B) {
	logger := phuslu.Logger{Level: phuslu.InfoLevel, Writer: phuslu.IOWriter{Writer: io.Discard}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("rate", "15").Int("low", 16).Float32("high", 123.2).Msg(msg)
	}
}

func BenchmarkLogrusPositive(b *testing.B) {
	logger := logrus.New()
	logger.Out = io.Discard
	logger.Level = logrus.InfoLevel

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkGookitSlogPositive(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.InfoLevel}),
	)
	logger.ReportCaller = false

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}
