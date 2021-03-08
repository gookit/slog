# slog

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/slog?style=flat-square)
[![GoDoc](https://godoc.org/github.com/gookit/slog?status.svg)](https://pkg.go.dev/github.com/gookit/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/slog)](https://goreportcard.com/report/github.com/gookit/slog)
[![Unit-Tests](https://github.com/gookit/slog/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/slog/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/slog)](https://github.com/gookit/slog)

📑 Lightweight, extensible, configurable logging library written in Go

## [中文说明](README.zh-CN.md)

中文说明请阅读 [README.zh-CN](README.zh-CN.md)

## Features

- Simple, directly available without configuration
- Support common log level processing. eg: `trace` `debug` `info` `notice` `warn` `error` `fatal` `panic`
- Supports adding multiple `Handler` log processing at the same time, outputting logs to different places
- Support any extension of `Handler` `Formatter` as needed
- Support to custom log messages `Handler`
- Support to custom log message `Formatter`
  - Built-in `json` `text` two log record formatting `Formatter`
- Has built-in common log write processing program
  - `console` output logs to the console, supports color output
  - `stream` output logs to the specified `io.Writer`
  - `simple_file` output logs to file, no buffer Write directly to file
  - `file` output logs to file. By default, `buffer` is enabled.
  - `size_rotate_file` output logs to file, and supports rotating files by size. By default, `buffer` is enabled.
  - `time_rotate_file` output logs to file, and supports rotating files by time. By default, `buffer` is enabled.
  - `rotate_file` output logs to file, and supports rotating files by time and size. By default, `buffer` is enabled.

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

- Change output format

Change the default logger output format.

```go
slog.GetFormatter().(*slog.TextFormatter).Template = slog.NamedTemplate
```

**Output:**

![](_example/images/console-color-log1.png)

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

## Logs to file

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func myfunc() {
	h1 := handler.MustFileHandler("/tmp/error.log", true)
	h1.Levels = slog.Levels{slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}

	h2 := handler.MustFileHandler("/tmp/info.log", true)
	h1.Levels = slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	// add logs
	slog.Info("info message")
	slog.Error("error message")
}
```

----------

## How to use handler

### Create handler

```go
h1, err := handler.NewSimpleFile("info.log")

h2, err := handler.NewFileHandler("error.log")

h3, err := handler.NewFileHandler("error.log")
```

### Push handler to logger

```go
	// append to logger
	l := slog.PushHandler(h)

	// logging messages
	slog.Info("info message")
	slog.Warn("warn message")
```

### New logger with handlers

```go
	// for new logger
	l := slog.NewWithHandlers(h1, h2)

	// logging messages
	l.Info("info message")
	l.Warn("warn message")
```


## Built-in Handlers

### BufferedHandler

`BufferedHandler` - can wrapper an `io.WriteCloser` as an `slog.Handler`

```go
package mypkg
import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func myfunc() {
	fpath := "./testdata/buffered-os-file.log"

	file, err := handler.QuickOpenFile(fpath)
	assert.NoError(t, err)

	bh := handler.NewBuffered(file, 2048)

	// new logger
	l := slog.NewWithHandlers(bh)

	// logging messages
	l.Info("buffered info message")
	l.Warn("buffered warn message")
}
```

### ConsoleHandler

`ConsoleHandler` - output logs to the console terminal. support color by [gookit/color](https://github.com/gookit/color).

Create:

```go
func NewConsoleHandler(levels []slog.Level) *ConsoleHandler
```

### EmailHandler

`NewEmailHandler` - output logs to email.

Create:

```go
h := handler.NewEmailHandler(from EmailOption, toAddresses []string)
```

### RotateFileHandler

`RotateFileHandler` - output log messages to file.

Create:

```go
func NewRotateFile(filepath string) (*RotateFileHandler, error)
func NewRotateFileHandler(filepath string) (*RotateFileHandler, error)
```

### TimeRotateFileHandler

`TimeRotateFileHandler` - output log messages to file.

Create:

```go
func NewTimeRotateFile(filepath string) (*TimeRotateFileHandler, error)
func NewTimeRotateFileHandler(filepath string) (*TimeRotateFileHandler, error)
```

The rotating files format support:

```go
const (
	EveryDay rotateTime = iota
	EveryHour
	Every30Minutes
	Every15Minutes
	EveryMinute
	EverySecond // only use for tests
)
```

file examples:

```text
time-rotate-file.log
time-rotate-file.log.20201229_155753
time-rotate-file.log.20201229_155754
```

### SizeRotateFileHandler

`SizeRotateFileHandler` - output log messages to file.

Create:

```go
func NewSizeRotateFile(filepath string) (*SizeRotateFileHandler, error)
func NewSizeRotateFileHandler(filepath string) (*SizeRotateFileHandler, error)
```

The rotating files format is `filename.log.yMD_0000N`. such as:

```text
size-rotate-file.log
size-rotate-file.log.122915_00001
size-rotate-file.log.122915_00002
```

### SimpleFileHandler

`SimpleFileHandler` - direct write log messages to a file. _Not recommended for production environment_

Create:

```go
func NewSimpleFile(filepath string) (*SimpleFileHandler, error)
func NewSimpleFileHandler(filepath string) (*SimpleFileHandler, error)
```

## Custom Logger

### Create New Logger

You can create a new instance of `slog.Logger`:

- Method 1：

```go
l := slog.New()
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- Method 2：

```go
l := slog.NewWithName("myLogger")
// add handlers ...
h1 := handler.NewConsoleHandler(slog.AllLevels)
l.AddHandlers(h1)
```

- Method 3：

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

### Create New Handler

you only need implement the `slog.Handler` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

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

### Create New Processor

you only need implement the `slog.Processor` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// AddHostname to record
func AddHostname() slog.Processor {
	hostname, _ := os.Hostname()

	return slog.ProcessorFunc(func(record *slog.Record) {
		record.AddField("hostname", hostname)
	})
}
```

add the processor:

```go
slog.AddProcessor(mypkg.AddHostname())
```

or:

```go
l := slog.New()
l.AddProcessor(mypkg.AddHostname())
```

### Create New Formatter

you only need implement the `slog.Formatter` interface:

```go
package mypkg

import (
	"github.com/gookit/slog"
)

type MyFormatter struct {
}

func (f *Formatter) Format(r *slog.Record) error {
	// format Record to text/JSON or other format.
}
```

add the formatter:

```go
slog.SetFormatter(&mypkg.MyFormatter{})
```

OR:

```go
l := slog.New()
h := &MyHandler{}
h.SetFormatter(&mypkg.MyFormatter{})

l.AddHander(h)
```

----------

## Introduction

slog handle workflow(like monolog):

```text
         Processors
Logger -{
         Handlers -{ With Formatters
```

### Processor

`Processor` - Logging `Record` processor.

You can use it to perform additional operations on the record before the log record reaches the Handler for processing, such as adding fields, adding extended information, etc.

`Processor` definition:

```go
// Processor interface definition
type Processor interface {
	// Process record
	Process(record *Record)
}
```

`Processor` func wrapper definition:

```go
// ProcessorFunc wrapper definition
type ProcessorFunc func(record *Record)

// Process record
func (fn ProcessorFunc) Process(record *Record) {
	fn(record)
}
```

Here we use the built-in processor `slog.AddHostname` as an example, it can add a new field `hostname` to each log record.

```go
slog.AddProcessor(slog.AddHostname())

slog.Info("message")
```

**Output:**

```json
{"channel":"application","level":"INFO","datetime":"2020/07/17 12:01:35","hostname":"InhereMac","data":{},"extra":{},"message":"message"}
```

### Handler

`Handler` - Log processor, each log will be processed by `Handler.Handle()`, where you can send the log to the console, file, or remote server.

> You can customize any `Handler` you want, you only need to implement the `slog.Handler` interface.

```go
// Handler interface definition
type Handler interface {
	// Close handler
	io.Closer
	// Flush logs to disk
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// all records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
}
```

> Note: Remember to add the `Handler` to the logger instance before the log records will be processed by the `Handler`.

### Formatter

`Formatter` - Log data formatting.

It is usually set in `Handler`, which can be used to format log records, convert records into text, JSON, etc., `Handler` then writes the formatted data to the specified place.

`Formatter` definition:

```go
// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}
```

`Formatter` function wrapper:

```go
// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format an record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}
```

**JSON Formatter**

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

**Text Formatter**

default templates:

```go
const DefaultTemplate = "[{{datetime}}] [{{channel}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n"
const NamedTemplate = "{{datetime}} channel={{channel}} level={{level}} [file={{caller}}] message={{message}} data={{data}}\n"
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
