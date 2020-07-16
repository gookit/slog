# slog

A simple log library for Go.

> Inspired the projects [Seldaek/monolog](https://github.com/Seldaek/monolog) and [sirupsen/logrus](https://github.com/sirupsen/logrus). Thank you very much

## [中文说明](README.zh-CN.md)

中文说明请阅读 [README.zh-CN](README.zh-CN.md)

## Features

- Simple, directly available without configuration
- Multiple `Handler` log handlers can be added at the same time to output logs to different places
- You can arbitrarily extend the `Handler` `Formatter` you need
- Support to support custom `Handler`
- Support to support custom `Formatter`

## GoDoc

- [godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## Install

```bash
go get github.com/gookit/slog
```

## Usage

`slog` is very simple to use and can be used without any configuration

## Quick Start

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Infof("info log %s", "message")
	slog.Debugf("debug %s", "message")
}
```

**Output:**

```text
[2020/07/16 12:19:33] [application] [INFO] info log message  
[2020/07/16 12:19:33] [application] [WARNING] warning log message  
[2020/07/16 12:19:33] [application] [INFO] info log message  
[2020/07/16 12:19:33] [application] [DEBUG] debug message  
```

### Use JSON Format

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	// use JSON formatter
	slog.SetFormatter(slog.NewJSONFormatter())

	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.WithData(slog.M{
		"key0": 134,
		"key1": "abc",
	}).Infof("info log %s", "message")

	r := slog.WithFields(slog.M{
		"category": "service",
		"IP": "127.0.0.1",
	})
	r.Infof("info %s", "message")
	r.Debugf("debug %s", "message")
}
```

**Output:**

```text
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"WARNING","message":"warning log message"}
{"channel":"application","data":{"key0":134,"key1":"abc"},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"DEBUG","message":"debug message"}
```

## Workflow

```text
         Processors
Logger -{
         Handlers -{ With Formatters
```

## Processor

## Handler

## Formatter

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
