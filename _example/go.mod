module slog_example

go 1.17

require (
	github.com/golang/glog v1.0.0
	github.com/gookit/goutil v0.5.2
	github.com/gookit/slog v0.3.2
	github.com/phuslu/log v1.0.67
	github.com/rs/zerolog v1.26.1
	github.com/sirupsen/logrus v1.8.1
	github.com/syyongx/llog v0.0.0-20200222114215-e8f9f86ac0a3
	go.uber.org/zap v1.21.0
)

require (
	github.com/gookit/color v1.5.0 // indirect
	github.com/gookit/gsr v0.0.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace github.com/gookit/slog => ../
