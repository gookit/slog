module slog_example

go 1.14

require (
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529
	github.com/gookit/slog v0.1.5
	github.com/kr/pretty v0.1.0 // indirect
	github.com/phuslu/log v1.0.67
	github.com/rs/zerolog v1.22.0
	github.com/sirupsen/logrus v1.8.1
	github.com/syyongx/llog v0.0.0-20200222114215-e8f9f86ac0a3
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/gookit/slog => ../
