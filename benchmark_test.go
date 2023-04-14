package slog_test

import (
	"io"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem bench_test.go
//
// code refer:
//
//	https://github.com/phuslu/log
var msg = "The quick brown fox jumps over the lazy dog"

func BenchmarkGookitSlogNegative(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.ErrorLevel}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func TestLogger_Info_Negative(t *testing.T) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, []slog.Level{slog.ErrorLevel}),
	)

	logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
}

func BenchmarkGookitSlogPositive(b *testing.B) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, slog.NormalLevels),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}
}

func BenchmarkTextFormatter_Format(b *testing.B) {
	r := newLogRecord("TEST_LOG_MESSAGE")
	f := slog.NewTextFormatter()
	// 1284 ns/op  456 B/op          11 allocs/op
	// On use DefaultTemplate

	// 304.4 ns/op   200 B/op           2 allocs/op
	// f.SetTemplate("{{datetime}} {{message}}")

	// 271.3 ns/op  200 B/op           2 allocs/op
	// f.SetTemplate("{{datetime}}")
	// f.SetTemplate("{{message}}")
	dump.P(f.Template())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := f.Format(r)
		if err != nil {
			panic(err)
		}
	}
}

func TestLogger_Info_Positive(t *testing.T) {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(io.Discard, slog.NormalLevels),
	)

	logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
}
