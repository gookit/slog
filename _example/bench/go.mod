module slog_bench

go 1.12

require (
	github.com/gookit/slog v0.1.3
	github.com/phuslu/log v1.0.67
	github.com/rs/zerolog v1.20.0
	go.uber.org/zap v1.16.0
)


replace (
	github.com/gookit/slog => ../../
)
