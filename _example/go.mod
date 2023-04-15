module slog_example

go 1.18

require (
	github.com/golang/glog v1.1.1
	github.com/gookit/goutil v0.6.8
	github.com/gookit/slog v0.5.0
	github.com/phuslu/log v1.0.67
	github.com/rs/zerolog v1.29.0
	github.com/sirupsen/logrus v1.9.0
	github.com/syyongx/llog v0.0.0-20200222114215-e8f9f86ac0a3
	go.uber.org/zap v1.24.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/gookit/color v1.5.3 // indirect
	github.com/gookit/gsr v0.0.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)

replace github.com/gookit/slog => ../
