# slog

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/slog?style=flat-square)
[![GoDoc](https://godoc.org/github.com/gookit/slog?status.svg)](https://pkg.go.dev/github.com/gookit/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/slog)](https://goreportcard.com/report/github.com/gookit/slog)
[![Unit-Tests](https://github.com/gookit/slog/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/slog/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/slog)](https://github.com/gookit/slog)
[![Coverage Status](https://coveralls.io/repos/github/gookit/slog/badge.svg?branch=master)](https://coveralls.io/github/gookit/slog?branch=master)

📑 Lightweight, extensible, configurable logging library written in Golang.

**Output in console:**

![console-log-all-level](_example/images/console-log-all-level.png)

## Features

- Simple, directly available without configuration
- Support common log level processing.
  - eg: `trace` `debug` `info` `notice` `warn` `error` `fatal` `panic`
- Support any extension of `Handler` `Formatter` as needed
- Supports adding multiple `Handler` log processing at the same time, outputting logs to different places
- Support to custom log message `Formatter`
  - Built-in `json` `text` two log record formatting `Formatter`
- Support to custom build log messages `Handler`
  - The built-in `handler.Config` `handler.Builder` can easily and quickly build the desired log handler
- Has built-in common log write handler program
  - `console` output logs to the console, supports color output
  - `writer` output logs to the specified `io.Writer`
  - `file` output log to the specified file, optionally enable `buffer` to buffer writes
  - `simple` output log to the specified file, write directly to the file without buffering
  - `rotate_file` outputs logs to the specified file, and supports splitting files by time and size, and `buffer` buffered writing is enabled by default
  - See [./handler](./handler) folder for more built-in implementations
- Benchmark performance test please see [Benchmarks](#benchmarks)

### Output logs to file

- Support enabling `buffer` for log writing
- Support splitting log files by `time` and `size`
- Support configuration to compress log files via `gzip`
- Support clean old log files by `BackupNum` `BackupTime`

### `rotatefile` subpackage

- The `rotatefile` subpackage is a stand-alone tool library with file splitting, cleaning, and compressing backups
- `rotatefile.Writer` can also be directly wrapped and used in other logging libraries. For example: `log`, `glog`, `zap`, etc.
- `rotatefile.FilesClear` is an independent file cleaning backup tool, which can be used in other places (such as other program log cleaning such as PHP)
- For more usage, please see [rotatefile](rotatefile/README.md)

## [中文说明](README.zh-CN.md)

中文说明请阅读 [README.zh-CN](README.zh-CN.md)

## GoDoc

- [Godoc for github](https://pkg.go.dev/github.com/gookit/slog?tab=doc)

## Install

```bash
go get github.com/gookit/slog
```

## Quick Start

`slog` is very simple to use and can be used without any configuration

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
[2020/07/16 12:19:33] [application] [INFO] [main.go:7] info log message  
[2020/07/16 12:19:33] [application] [WARNING] [main.go:8] warning log message  
[2020/07/16 12:19:33] [application] [INFO] [main.go:9] info log message  
[2020/07/16 12:19:33] [application] [DEBUG] [main.go:10] debug message  
```

### Console Color

You can enable color on output logs to console. _This is default_

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

**Output:**

![](_example/images/console-color-log.png)

### Change log output style

Above is the `Formatter` setting that changed the default logger.

> You can also create your own logger and append `ConsoleHandler` to support printing logs to the console:

```go
h := handler.NewConsoleHandler(slog.AllLevels)
l := slog.NewWithHandlers(h)

l.Trace("this is a simple log message")
l.Debug("this is a simple log message")
```

Change the default logger log output style:

```go
h.Formatter().(*slog.TextFormatter).SetTemplate(slog.NamedTemplate)
```

**Output:**

![](_example/images/console-color-log1.png)

> Note: `slog.TextFormatter` uses a template string to format the output log, so the new field output needs to adjust the template at the same time.

### Use JSON Format

`slog` also has a built-in `Formatter` for JSON format. If not specified, the default is to use `TextFormatter` to format log records.

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

- `Logger` - log dispatcher. One logger can register multiple `Handler`, `Processor`
- `Record` - log records, each log is a `Record` instance.
- `Processor` - enables extended processing of log records. It is called before the log `Record` is processed by the `Handler`.
  - You can use it to perform additional operations on `Record`, such as: adding fields, adding extended information, etc.
- `Handler` - log handler, each log will be processed by `Handler.Handle()`.
  - Here you can send logs to console, file, remote server, etc.
- `Formatter` - logging data formatting process.
  - Usually set in `Handler`, it can be used to format log records, convert records into text, JSON, etc., `Handler` then writes the formatted data to the specified place.
  - `Formatter` is not required. You can do without it and handle logging directly in `Handler.Handle()`.

**Simple structure of log scheduler**：

```text
          Processors
Logger --{
          Handlers --|- Handler0 With Formatter0
                     |- Handler1 With Formatter1
                     |- Handler2 (can also without Formatter)
                     |- ... more
```

> Note: Be sure to remember to add `Handler`, `Processor` to the logger instance and log records will be processed by `Handler`.

### Processor

`Processor` interface:

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

> You can use it to perform additional operations on the Record before the log `Record` reaches the `Handler` for processing, such as: adding fields, adding extended information, etc.

Add processor to logger:

```go
slog.AddProcessor(slog.AddHostname())

// or
l := slog.New()
l.AddProcessor(slog.AddHostname())
```

The built-in processor `slog.AddHostname` is used here as an example, which can add a new field `hostname` on each log record.

```go
slog.AddProcessor(slog.AddHostname())
slog.Info("message")
```

Output, including new fields `"hostname":"InhereMac"`：

```json
{"channel":"application","level":"INFO","datetime":"2020/07/17 12:01:35","hostname":"InhereMac","data":{},"extra":{},"message":"message"}
```

### Handler

`Handler` interface:

> You can customize any `Handler` you want, just implement the `slog.Handler` interface.

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

### Formatter

`Formatter` interface:

```go
// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}
```

Function wrapper type：

```go
// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format a log record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}
```

**JSON formatter**

```go
type JSONFormatter struct {
	// Fields exported log fields.
	Fields []string
	// Aliases for output fields. you can change export field name.
	// item: `"field" : "output name"`
	// eg: {"message": "msg"} export field will display "msg"
	Aliases StringMap
	// PrettyPrint will indent all json logs
	PrettyPrint bool
	// TimeFormat the time format layout. default is time.RFC3339
	TimeFormat string
}
```

**Text formatter**

Default templates:

```go
const DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n"
const NamedTemplate = "{{datetime}} channel={{channel}} level={{level}} [file={{caller}}] message={{message}} data={{data}}\n"
```

Change template:

```go
myTemplate := "[{{datetime}}] [{{level}}] {{message}}"

f := slog.NewTextFormatter()
f.SetTemplate(myTemplate)
```

## Custom logger

Custom `Processor` and `Formatter` are relatively simple, just implement a corresponding method.

### Create new logger

`slog.Info, slog.Warn` and other methods use the default logger and output logs to the console by default.

You can create a brand-new instance of `slog.Logger`:

**Method 1**：

```go
l := slog.New()
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

**Method 2**：

```go
l := slog.NewWithName("myLogger")
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

**Method 3**：

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

### Create custom Handler

You only need to implement the `slog.Handler` interface to create a custom `Handler`.

You can quickly assemble your own Handler through the built-in `handler.LevelsWithFormatter` `handler.LevelWithFormatter` and other fragments of slog.

Examples:

> Use `handler.LevelsWithFormatter`, only need to implement `Close, Flush, Handle` methods

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

Add `Handler` to the logger to use:

```go
// add to default logger
slog.AddHander(&MyHandler{})

// or, add to custom logger:
l := slog.New()
l.AddHander(&MyHandler{})
```

## Use the built-in handlers

[./handler](handler) package has built-in common log handlers, which can basically meet most scenarios.

```go
// Output logs to console, allow render color.
func NewConsoleHandler(levels []slog.Level) *ConsoleHandler
// Send logs to email
func NewEmailHandler(from EmailOption, toAddresses []string) *EmailHandler
// Send logs to syslog
func NewSysLogHandler(priority syslog.Priority, tag string) (*SysLogHandler, error)
// A simple handler implementation that outputs logs to a given io.Writer
func NewSimpleHandler(out io.Writer, level slog.Level) *SimpleHandler
```

**Output log to file**:

```go
// Output log to the specified file, without buffering by default
func NewFileHandler(logfile string, fns ...ConfigFn) (h *SyncCloseHandler, err error)
// Output logs to the specified file in JSON format, without buffering by default
func JSONFileHandler(logfile string, fns ...ConfigFn) (*SyncCloseHandler, error)
// Buffered output log to specified file
func NewBuffFileHandler(logfile string, buffSize int, fns ...ConfigFn) (*SyncCloseHandler, error)
```

> TIP: `NewFileHandler` `JSONFileHandler` can also enable write buffering by passing in fns `handler.WithBuffSize(buffSize)`

**Output log to file and rotate automatically**:

```go
// Automatic rotating according to file size
func NewSizeRotateFile(logfile string, maxSize int, fns ...ConfigFn) (*SyncCloseHandler, error)
// Automatic rotating according to time
func NewTimeRotateFile(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error)
// It supports configuration to rotate according to size and time. 
// The default setting file size is 20M, and the default automatic splitting time is 1 hour (EveryHour).
func NewRotateFileHandler(logfile string, rt rotatefile.RotateTime, fns ...ConfigFn) (*SyncCloseHandler, error)
```

> TIP: By passing in `fns ...ConfigFn`, more options can be set, such as log file retention time, log write buffer size, etc. For detailed settings, see the `handler.Config` structure

### Logs to file

Output log to the specified file, `buffer` buffered writing is not enabled by default. Buffering can also be enabled by passing in a parameter.

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func main() {
	defer slog.MustClose()

	// DangerLevels contains: slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel
	h1 := handler.MustFileHandler("/tmp/error.log", handler.WithLogLevels(slog.DangerLevels))
	// custom log format
	// f := h1.Formatter().(*slog.TextFormatter)
	f := slog.AsTextFormatter(h1.Formatter())
	f.SetTemplate("your template format\n")

	// NormalLevels contains: slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel
	h2 := handler.MustFileHandler("/tmp/info.log", handler.WithLogLevels(slog.NormalLevels))

	// register handlers
	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}
```

> **Note**: If write buffering `buffer` is enabled, be sure to call `logger.Close()` at the end of the program to flush the contents of the buffer to the file.

### Log to file with automatic rotating

`slog/handler` also has a built-in output log to a specified file, and supports splitting files by time and size at the same time.
By default, `buffer` buffered writing is enabled

```go
func Example_rotateFileHandler() {
	h1 := handler.MustRotateFile("/tmp/error.log", handler.EveryHour, handler.WithLogLevels(slog.DangerLevels))
	h2 := handler.MustRotateFile("/tmp/info.log", handler.EveryHour, handler.WithLogLevels(slog.NormalLevels))

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}
```

Example of file name sliced by time:

```text
time-rotate-file.log
time-rotate-file.log.20201229_155753
time-rotate-file.log.20201229_155754
```

Example of a filename cut by size, in the format `filename.log.HIS_000N`. For example:

```text
size-rotate-file.log
size-rotate-file.log.122915_0001
size-rotate-file.log.122915_0002
```

### Use rotatefile on another logger

`rotatefile.Writer` can also be used with other logging packages, such as: `log`, `glog`, etc. 

For example, using `rotatefile` on golang `log`:

```go
package main

import (
  "log"

  "github.com/gookit/slog/rotatefile"
)

func main() {
	logFile := "testdata/go_logger.log"
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err) 
	}

	log.SetOutput(writer)
	log.Println("log message")
}
```

### Quickly create a Handler based on config

This is config struct for create a Handler:

```go
// Config struct
type Config struct {
	// Logfile for write logs
	Logfile string `json:"logfile" yaml:"logfile"`
	// LevelMode for filter log record. default LevelModeList
	LevelMode uint8 `json:"level_mode" yaml:"level_mode"`
	// Level value. use on LevelMode = LevelModeValue
	Level slog.Level `json:"level" yaml:"level"`
	// Levels for log record
	Levels []slog.Level `json:"levels" yaml:"levels"`
	// UseJSON for format logs
	UseJSON bool `json:"use_json" yaml:"use_json"`
	// BuffMode type name. allow: line, bite
	BuffMode string `json:"buff_mode" yaml:"buff_mode"`
	// BuffSize for enable buffer, unit is bytes. set 0 to disable buffer
	BuffSize int `json:"buff_size" yaml:"buff_size"`
	// RotateTime for rotate file, unit is seconds.
	RotateTime rotatefile.RotateTime `json:"rotate_time" yaml:"rotate_time"`
	// MaxSize on rotate file by size, unit is bytes.
	MaxSize uint64 `json:"max_size" yaml:"max_size"`
	// Compress determines if the rotated log files should be compressed using gzip.
	// The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`
	// BackupNum max number for keep old files.
	// 0 is not limit, default is 20.
	BackupNum uint `json:"backup_num" yaml:"backup_num"`
	// BackupTime max time for keep old files. unit is hours
	// 0 is not limit, default is a week.
	BackupTime uint `json:"backup_time" yaml:"backup_time"`
	// RenameFunc build filename for rotate file
	RenameFunc func(filepath string, rotateNum uint) string
}
```

**Examples**:

```go
	testFile := "testdata/error.log"

	h := handler.NewEmptyConfig(
			handler.WithLogfile(testFile),
			handler.WithBuffSize(1024*8),
			handler.WithRotateTimeString("1hour"),
			handler.WithLogLevels(slog.DangerLevels),
		).
		CreateHandler()

	l := slog.NewWithHandlers(h)
```

**About BuffMode**

`Config.BuffMode` The name of the BuffMode type to use. Allow: line, bite

- `BuffModeBite`: Buffer by bytes, when the number of bytes in the buffer reaches the specified size, write the contents of the buffer to the file
- `BuffModeLine`: Buffer by line, when the buffer size is reached, always ensure that a complete line of log content is written to the file (to avoid log content being truncated)

### Use Builder to quickly create Handler

Use `handler.Builder` to easily and quickly create Handler instances.

```go
	testFile := "testdata/info.log"

	h := handler.NewBuilder().
		WithLogfile(testFile).
		WithLogLevels(slog.NormalLevels).
		WithBuffSize(1024*8).
		WithRotateTime(rotatefile.Every30Min).
		WithCompress(true).
		Build()

	l := slog.NewWithHandlers(h)
```

## Extension packages

Package `bufwrite`:

- `bufwrite.BufIOWriter` additionally implements `Sync(), Close()` methods by wrapping go's `bufio.Writer`, which is convenient to use
- `bufwrite.LineWriter` refer to the implementation of `bufio.Writer` in go, which can support flushing the buffer by line, which is more useful for writing log files

Package `rotatefile`:

- `rotatefile.Writer` implements automatic cutting of log files according to size and specified time, and also supports automatic cleaning of log files
  - `handler/rotate_file` is to use it to cut the log file

### Use rotatefile on other log package

Of course, the rotatefile.Writer can be use on other log package, such as: `log`, `glog` and more.

Examples, use rotatefile on golang `log`:

```go
package main

import (
  "log"

  "github.com/gookit/slog/rotatefile"
)

func main() {
	logFile := "testdata/another_logger.log"
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err) 
	}

	log.SetOutput(writer)
	log.Println("log message")
}
```

## Testing and benchmark

### Unit tests

run unit tests:

```bash
go test ./...
```

### Benchmarks

Benchmark code at [_example/bench_loglibs_test.go](_example/bench_loglibs_test.go)

```bash
make test-bench
```

Benchmarks for `slog` and other log packages:

> **Note**: test and record ad 2023.04.13

```shell
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-3740QM CPU @ 2.70GHz
BenchmarkZapNegative
BenchmarkZapNegative-4                   8381674              1429 ns/op             216 B/op          3 allocs/op
BenchmarkZapSugarNegative
BenchmarkZapSugarNegative-4              8655980              1383 ns/op             104 B/op          4 allocs/op
BenchmarkZeroLogNegative
BenchmarkZeroLogNegative-4              14173719               849.8 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogNegative
BenchmarkPhusLogNegative-4              27456256               451.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusNegative
BenchmarkLogrusNegative-4                2550771              4784 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogNegative
>>>> BenchmarkGookitSlogNegative-4       8798220              1375 ns/op             120 B/op          3 allocs/op
BenchmarkZapPositive
BenchmarkZapPositive-4                  10302483              1167 ns/op             192 B/op          1 allocs/op
BenchmarkZapSugarPositive
BenchmarkZapSugarPositive-4              3833311              3154 ns/op             344 B/op          7 allocs/op
BenchmarkZeroLogPositive
BenchmarkZeroLogPositive-4              14120524               846.7 ns/op             0 B/op          0 allocs/op
BenchmarkPhusLogPositive
BenchmarkPhusLogPositive-4              27152686               434.9 ns/op             0 B/op          0 allocs/op
BenchmarkLogrusPositive
BenchmarkLogrusPositive-4                2601892              4691 ns/op             608 B/op         17 allocs/op
BenchmarkGookitSlogPositive
>>>> BenchmarkGookitSlogPositive-4            8997104              1340 ns/op             120 B/op          3 allocs/op
PASS
ok      command-line-arguments  167.095s
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

The projects is heavily inspired by follow packages:

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
