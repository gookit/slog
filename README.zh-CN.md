# slog

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/slog?style=flat-square)
[![GoDoc](https://godoc.org/github.com/gookit/slog?status.svg)](https://pkg.go.dev/github.com/gookit/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/slog)](https://goreportcard.com/report/github.com/gookit/slog)
[![Unit-Tests](https://github.com/gookit/slog/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/slog/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/slog)](https://github.com/gookit/slog)
[![Coverage Status](https://coveralls.io/repos/github/gookit/slog/badge.svg?branch=master)](https://coveralls.io/github/gookit/slog?branch=master)

📑 Go 实现的一个易于使用的，易扩展、可配置的日志库

**控制台日志效果:**

![console-log-all-level](_example/images/console-log-all-level.png)

## 功能特色

- 简单，无需配置，开箱即用
- 支持常用的日志级别处理
  - 如： `trace` `debug` `info` `notice` `warn` `error` `fatal` `panic`
- 可以任意扩展自己需要的 `Handler` `Formatter` 
- 支持同时添加多个 `Handler` 日志处理，输出日志到不同的地方
- 支持自定义构建 `Handler` 处理器
  - 内置的 `handler.Config` `handler.Builder`,可以方便快捷的构建想要的日志处理器
- 支持自定义 `Formatter` 格式化处理
  - 内置了 `json` `text` 两个日志记录格式化 `Formatter`
- 已经内置了常用的日志处理器
  - `console` 输出日志到控制台，支持色彩输出
  - `writer` 输出日志到指定的 `io.Writer`
  - `simple_file` 输出日志到指定文件，无缓冲直接写入文件
  - `file` 输出日志到指定文件，默认启用 `buffer` 缓冲写入
  - `rotate_file` 输出日志到指定文件，并且同时支持按时间、按大小分割文件，默认启用 `buffer` 缓冲写入

**扩展工具包**

`bufwrite` 包:

- `bufwrite.BufIOWriter` 通过包装go的 `bufio.Writer` 额外实现了 `Sync(), Close()` 方法，方便使用
- `bufwrite.LineWriter` 参考go的 `bufio.Writer` 实现, 可以支持按行刷出缓冲，对于写日志文件更有用

`rotatefile` 包:

- `rotatefile.Writer` 实现对日志文件按大小和指定时间进行自动切割，同时也支持自动清理日志文件
  - `RotateHandler` 即是通过使用它对日志文件进行切割处理

> NEW: `v0.3.0` 废弃原来实现的纷乱的各种handler,统一抽象为
> `FlushCloseHandler` `SyncCloseHandler` `WriteCloserHandler` `IOWriterHandler` 
> 几个支持不同类型writer的处理器。让构建自定义 Handler 更加简单，内置的handlers也基本上由它们组成。

## [English](README.md)

English instructions please read [README](README.md)

## GoDoc

- [Godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## 安装

```bash
go get github.com/gookit/slog
```

## 快速开始

`slog` 使用非常简单，无需任何配置即可使用。

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

### 更改日志输出样式

上面是更改了默认logger的 `Formatter` 设置。

你也可以创建自己的logger，并追加 `ConsoleHandler` 来支持打印日志到控制台：

```go
h := handler.NewConsoleHandler(slog.AllLevels)
l := slog.NewWithHandlers()

l.Trace("this is a simple log message")
l.Debug("this is a simple log message")
```

更改默认的logger日志输出样式:

```go
h.GetFormatter().(*slog.TextFormatter).Template = slog.NamedTemplate
```

**输出预览:**

![](_example/images/console-color-log1.png)

> 注意： `slog.TextFormatter` 使用模板字符串来格式化输出日志，因此新增字段输出需要同时调整模板。

### 使用JSON格式

slog 也内置了 JSON 格式的 `Formatter`。若不特别指定，默认都是使用 `TextFormatter` 格式化日志记录。

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

## 架构说明

- `Logger` - 日志调度器. 一个logger可以注册多个 `Handler`,`Processor`
- `Record` - 日志记录，每条日志就是一个 `Record` 实例。
- `Processor` - 可以对日志记录进行扩展处理。它在日志 `Record` 被 `Handler` 处理之前调用。
  - 你可以使用它对 `Record` 进行额外的操作，比如：新增字段，添加扩展信息等
- `Handler` - 日志处理器，每条日志都会经过 `Handler.Handle()` 处理。
  - 在这里你可以将日志发送到 控制台，文件，远程服务器等等。
- `Formatter` - 日志记录数据格式化处理。
  - 通常设置于 `Handler` 中，可以用于格式化日志记录，将记录转成文本，JSON等，`Handler` 再将格式化后的数据写入到指定的地方。
  - `Formatter` 不是必须的。你可以不使用它,直接在 `Handler.Handle()` 中对日志记录进行处理。

**日志调度器简易结构**：

```text
          Processors
Logger --{
          Handlers --{ With Formatter
```

> 注意：一定要记得将 `Handler`, `Processor` 添加注册到 logger 实例上，日志记录才会经过 `Handler` 处理。

### Processor 定义

`Processor` 接口定义如下:

```go
// Processor interface definition
type Processor interface {
	// Process record
	Process(record *Record)
}

// ProcessorFunc definition
type ProcessorFunc func(record *Record)

// Process record
func (fn ProcessorFunc) Process(record *Record) {
	fn(record)
}
```

> 你可以使用它在日志 `Record` 到达 `Handler` 处理之前，对Record进行额外的操作，比如：新增字段，添加扩展信息等

这里使用内置的processor `slog.AddHostname` 作为示例，它可以在每条日志记录上添加新字段 `hostname`。

```go
slog.AddProcessor(slog.AddHostname())
slog.Info("message")
```

输出效果,包含新增字段 `"hostname":"InhereMac"`：

```json
{"channel":"application","level":"INFO","datetime":"2020/07/17 12:01:35","hostname":"InhereMac","data":{},"extra":{},"message":"message"}
```

### Handler 定义

`Handler` 接口定义如下:

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

### Formatter 定义

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

// Format a log record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}
```

## 自定义日志

自定义 Processor 和 自定义 Formatter 都比较简单，实现一个对应方法即可。

### 创建自定义 Logger实例

`slog.Info` 等方法，使用的默认logger，并且默认输出日志到控制台。 

你可以创建一个全新的 `slog.Logger` 实例：

**方式1**：

```go
l := slog.New()
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

**方式2**：

```go
l := slog.NewWithName("myLogger")
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

**方式3**：

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

你只需要实现 `slog.Handler` 接口即可创建自定义 `Handler`。你可以通过 slog内置的
`handler.LevelsWithFormatter` `handler.LevelWithFormatter`等片段快速的组装自己的 Handler。

示例:

> 使用了 `handler.LevelsWithFormatter`， 只还需要实现 `Close, Flush, Handle` 方法即可

```go
type MyHandler struct {
	handler.LevelsWithFormatter
    Output io.Writer
}

func (h *MyHandler) Handle(r *slog.Record) error {
	// you can write log message to file or send to remote.
}

func (h *MyHandler) Flush() error {}
func (h *MyHandler) Close() error {}
```

将 `Handler` 添加到 logger即可使用:

```go
// 添加到默认 logger
slog.AddHander(&MyHandler{})

// 或者添加到自定义 logger:
l := slog.New()
l.AddHander(&MyHandler{})
```

## 使用内置处理器

[./handler](handler) 包已经内置了常用的日志 Handler，基本上可以满足绝大部分场景。

```go
NewConsoleHandler(levels []slog.Level) // 输出日志到控制台
NewEmailHandler(from EmailOption, toAddresses []string) // 发送日志到email邮箱
NewSysLogHandler(priority syslog.Priority, tag string) (*SysLogHandler, error) // 发送日志到系统的syslog
NewSimpleHandler(out io.Writer, level slog.Level) *SimpleHandler // 一个简单的handler实现，输出日志到给定的 io.Writer
```

输出日志到文件:

```go
NewFileHandler(logfile string, fns ...ConfigFn) (h *SyncCloseHandler, err error)  // 输出日志到指定文件，默认不带缓冲
JSONFileHandler(logfile string, fns ...ConfigFn) (*SyncCloseHandler, error)  // 输出日志到指定文件且格式为JSON，默认不带缓冲
NewBuffFileHandler(logfile string, buffSize int, fns ...ConfigFn) (*SyncCloseHandler, error)  // 带缓冲的输出日志到指定文件
```

> TIP: `NewFileHandler` `JSONFileHandler` 也可以通过传入 fns `handler.WithBuffSize(buffSize)` 启用写入缓冲

输出日志到文件并自动切割:

```go
NewSizeRotateFile(logfile string, maxSize int, fns ...ConfigFn) (*SyncCloseHandler, error) // 根据文件大小进行自动切割
NewTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error) // 根据时间进行自动切割
// 同时支持配置根据大小和时间进行切割, 默认设置文件大小是 20M，默认自动分割时间是 1小时(EveryHour)。
NewRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error)
```

> TIP: 通过传入 `fns ...ConfigFn` 可以设置更多选项，比如 日志文件保留时间, 日志写入缓冲大小等。 详细设置请看 `handler.Config` 结构体

### 输出日志到文件

`FileHandler` 输出日志到指定文件，默认不启用 `buffer` 缓冲写入。 也可以通过传入参数启用缓冲。

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func main() {
	defer slog.MustFlush()

	// DangerLevels 包含： slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel
	h1 := handler.MustFileHandler("/tmp/error.log", handler.WithLogLevels(slog.DangerLevels))

	// NormalLevels 包含： slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel
	h2 := handler.MustFileHandler("/tmp/info.log", handler.WithLogLevels(slog.NormalLevels))

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}
```

> 提示: 如果启用了写入缓冲 `buffer`，一定要在程序结束时调用 `logger.Flush()` 刷出缓冲区的内容到文件。

### 根据配置快速创建Handler实例



## 测试以及性能

### 单元测试

运行单元测试

```bash
go test -v ./...
```

### 性能压测

```bash
make test-bench
```

> record ad 2022.04.27

```text
% make test-bench
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                  128133166               93.97 ns/op          192 B/op          1 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              909583207               13.41 ns/op            0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              784099310               15.24 ns/op            0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4               289939296               41.60 ns/op           16 B/op          1 allocs/op
BenchmarkGookitSlogNegative
> BenchmarkGookitSlogNegative-4           29131203               417.4 ns/op           125 B/op          4 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                   9910075              1219 ns/op             192 B/op          1 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              13966810               871.0 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              26743148               446.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2658482              4481 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
> BenchmarkGookitSlogPositive-4            8349562              1441 ns/op             165 B/op          6 allocs/op
PASS
ok      command-line-arguments  146.669s
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

## Acknowledgment

实现参考了以下项目，非常感谢它们

- https://github.com/phuslu/log
- https://github.com/golang/glog
- https://github.com/sirupsen/logrus
- https://github.com/Seldaek/monolog
- https://github.com/syyongx/llog
- https://github.com/uber-go/zap
- https://github.com/rs/zerolog
- https://github.com/natefinch/lumberjack

## LICENSE

[MIT](LICENSE)
