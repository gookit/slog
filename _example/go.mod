module slog_example

go 1.19

require (
	github.com/golang/glog v1.2.5
	github.com/gookit/goutil v0.7.1
	github.com/gookit/slog v0.5.8
	github.com/phuslu/log v1.0.119
	github.com/rs/zerolog v1.34.0
	github.com/sirupsen/logrus v1.9.3
	github.com/syyongx/llog v0.0.0-20200222114215-e8f9f86ac0a3
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/gookit/color v1.6.0 // indirect
	github.com/gookit/gsr v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/gookit/slog => ../
