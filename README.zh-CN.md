# slog

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/slog?style=flat-square)
[![GoDoc](https://godoc.org/github.com/gookit/slog?status.svg)](https://pkg.go.dev/github.com/gookit/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/slog)](https://goreportcard.com/report/github.com/gookit/slog)
[![Unit-Tests](https://github.com/gookit/slog/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/slog/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/slog)](https://github.com/gookit/slog)

📑 Go 实现的开箱即用，易扩展、可配置的日志库

> 项目实现参考了 [Seldaek/monolog](https://github.com/Seldaek/monolog) and [sirupsen/logrus](https://github.com/sirupsen/logrus) ，非常感谢它们。

## [English](README.md)

English instructions please read [README](README.md)

## 功能特色

- 简单，无需配置，开箱即用
- 支持常用的日志级别处理。如： `trace` `debug` `info` `notice` `warn` `error` `fatal` `panic`
- 支持同时添加多个 `Handler` 日志处理，输出日志到不同的地方
- 可以任意扩展自己需要的 `Handler` `Formatter` 
- 支持自定义 `Handler` 处理程器
- 支持自定义 `Formatter` 格式化处理
  - 内置了 `json` `text` 两个日志记录格式化 `Formatter`
- 已经内置了常用的日志写入处理程序
  - `console` 输出日志到控制台，支持色彩输出
  - `stream` 输出日志到指定的 `io.Writer`
  - `simple_file` 输出日志到指定文件，无缓冲直接写入文件
  - `file` 输出日志到指定文件，默认启用 `buffer` 缓冲写入
  - `size_rotate_file` 输出日志到指定文件，并且支持按大小分割文件。默认启用 `buffer` 缓冲写入
  - `time_rotate_file` 输出日志到指定文件，并且支持按时间分割文件。默认启用 `buffer` 缓冲写入
  - `rotate_file` 输出日志到指定文件，并且同时支持按时间、按大小分割文件。默认启用 `buffer` 缓冲写入

## GoDoc

- [Godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## 安装

```bash
go get github.com/gookit/slog
```

## 使用

`slog` 使用非常简单，无需任何配置即可使用

## 快速开始

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

**输出预览:**

```text
[2020/07/16 12:19:33] [application] [INFO] info log message  
[2020/07/16 12:19:33] [application] [WARNING] warning log message  
[2020/07/16 12:19:33] [application] [INFO] info log message  
[2020/07/16 12:19:33] [application] [DEBUG] debug message  
```

### 启用控制台颜色

您可以在输出控制台日志时启用颜色输出，将会根据不同级别打印不同色彩。

```go
package main

import (
	"github.com/gookit/slog"
)

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
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

**输出预览:**

![](_example/images/console-color-log.png)

上面是更改了默认的 `Formatter` 设置。你也可以追加 `ConsoleHandler` 来支持打印日志到控制台：

```go
l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))

l.Trace("this is a simple log message")
l.Debug("this is a simple log message")
```

- 更改日志输出样式

更改默认的logger日志输出样式.

```go
slog.GetFormatter().(*slog.TextFormatter).Template = slog.NamedTemplate
```

**输出预览:**

![](_example/images/console-color-log1.png)

> 注意： `slog.TextFormatter` 使用模板字符串来格式化输出日志，因此新增字段输出需要同时调整模板

### 使用JSON格式

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

**输出预览:**

```text
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"channel":"application","data":{},"datetime":"2020/07/16 13:23:33","extra":{},"level":"WARNING","message":"warning log message"}
{"channel":"application","data":{"key0":134,"key1":"abc"},"datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info log message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"INFO","message":"info message"}
{"IP":"127.0.0.1","category":"service","channel":"application","datetime":"2020/07/16 13:23:33","extra":{},"level":"DEBUG","message":"debug message"}
```

## 自定义日志

## 创建自定义 Logger实例

你可以创建一个全新的 `slog.Logger` 实例：

- 方式1：

```go
l := slog.New()
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- 方式2：

```go
l := slog.NewWithName("myLogger")
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- 方式3：

```go
package main

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func main() {
	l := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	l.Info("message")
}
```

### 创建自定义 Handler

you only need implement the `slog.Handler` interface:

```go
type MyHandler struct {
	handler.LevelsWithFormatter
}

func (h *MyHandler) Handle(r *slog.Record) error {
	// you can write log message to file or send to remote.
}
```

add handler to default logger:

```go
slog.AddHander(&MyHandler{})
```

or add to custom logger:

```go
l := slog.New()
l.AddHander(&MyHandler{})
```

### 创建自定义 Processor

### 创建自定义 Formatter

## 架构说明

简易日志处理流程：

```text
         Processors
Logger -{
         Handlers -{ With Formatters
```

### Processor

`Processor` - 日志记录(`Record`)处理器。你可以使用它在日志 `Record` 到达 `Handler` 处理之前，对Record进行额外的操作，比如：新增字段，添加扩展信息等

这里使用内置的processor `slog.AddHostname` 作为示例，它可以在每条日志记录上添加新字段 `hostname`。

```go
slog.AddProcessor(slog.AddHostname())

slog.Info("message")
```

输出类似：

```json
{"channel":"application","level":"INFO","datetime":"2020/07/17 12:01:35","hostname":"InhereMac","data":{},"extra":{},"message":"message"}
```

### Handler

`Handler` - 日志处理器，每条日志都会经过 `Handler.Handle()` 处理，在这里你可以将日志发送到 控制台，文件，远程服务器。

> 你可以自定义任何想要的 `Handler`，只需要实现 `slog.Handler` 接口即可。

```go
// Handler interface definition
type Handler interface {
	io.Closer
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

`Formatter` 接口定义如下:

```go
// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}
```

函数包装类型：

```go
// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format an record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}
```

## Gookit packages

  - [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
  - [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
  - [gookit/gcli](https://github.com/gookit/gcli) Build CLI application, tool library, running CLI commands
  - [gookit/slog](https://github.com/gookit/slog) Lightweight, extensible, configurable logging library written in Go
  - [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
  - [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
  - [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
  - [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
  - [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
  - [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
  - [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
  - More, please see https://github.com/gookit

## Refer

- https://github.com/golang/glog
- https://github.com/sirupsen/logrus
- https://github.com/Seldaek/monolog
- https://github.com/syyongx/llog
- https://github.com/uber-go/zap
- https://github.com/rs/zerolog

## LICENSE

[MIT](LICENSE)
