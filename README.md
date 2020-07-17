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

### Enable Console Color

You can enable color on output logs to console.

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter().(*slog.TextFormatter)
		f.EnableColor = true
	})

	slog.Trace("this is a simple log message")
	slog.Debug("this is a simple log message")
	slog.Info("this is a simple log message")
	slog.Notice("this is a simple log message")
	slog.Warn("this is a simple log message")
	slog.Error("this is a simple log message")
	slog.Fatal("this is a simple log message")
}
```

**Output:**

![](_example/images/console-color-log.png)

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

## Introduction

slog handle workflow:

```text
         Processors
Logger -{
         Handlers -{ With Formatters
```

### Processor

`Processor` - 日志记录(`Record`)处理器。

你可以使用它在日志 `Record` 到达 `Handler` 处理之前，对Record进行额外的操作，比如：新增字段，添加扩展信息等

这里使用内置的processor `slog.AddHostname` 作为示例，它可以在每条日志记录上添加新字段 `hostname`。

```go
slog.AddProcessor(slog.AddHostname())

slog.Info("message")
```

输出类似：

```
{"channel":"application","data":{},"datetime":"2020/07/17 12:01:35","extra":{},"hostname":"InhereMac","level":"INFO","message":"message"}
```

### Handler

`Handler` - 日志处理器，每条日志都会经过 `Handler.Handle()` 处理，在这里你可以将日志发送到 控制台，文件，远程服务器。

> 你可以自定义任何想要的 `Handler`，只需要实现 `slog.Handler` 接口即可。

```go
// Handler interface definition
type Handler interface {
	// io.Closer
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// all records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
}
```

> 注意：一定要记得将 `Handler` 添加注册到 logger 实例上，日志记录才会经过 `Handler` 处理。

### Formatter

`Formatter` - 日志数据格式化。它通常设置于 `Handler` 中，可以用于格式化日志记录，将记录转成文本，JSON等，`Handler` 再将格式化后的数据写入到指定的地方。

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
