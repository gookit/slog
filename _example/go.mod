module slog_example

go 1.18

require (
	github.com/golang/glog v1.0.0
	github.com/gookit/goutil v0.6.4
	github.com/gookit/slog v0.4.0
	github.com/phuslu/log v1.0.67
	github.com/rs/zerolog v1.28.0
	github.com/sirupsen/logrus v1.9.0
	github.com/syyongx/llog v0.0.0-20200222114215-e8f9f86ac0a3
	go.uber.org/zap v1.23.0
)

require (
	github.com/gookit/color v1.5.2 // indirect
	github.com/gookit/gsr v0.0.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

replace github.com/gookit/slog => ../
