# slog

simple log

> Inspired the projects [Seldaek/monolog](https://github.com/Seldaek/monolog) and [sirupsen/logrus](https://github.com/sirupsen/logrus). Thank you very much

## GoDoc

- [godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## Install

```bash
go get github.com/gookit/slog
```

## Usage

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Infof("info log %s", "message")
}
```

## Workflow


```text
                     With Formatters
         Handlers -{ 
Logger -{
         Processors 
```

## Refer

- https://github.com/golang/glog
- https://github.com/Seldaek/monolog

## Related

- https://github.com/sirupsen/logrus
- https://github.com/uber-go/zap
- https://github.com/rs/zerolog
- https://github.com/syyongx/llog

## LICENSE

[MIT](LICENSE)
