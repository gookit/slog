package slog

import (
	"fmt"
	"io"
	"strings"
)

// M short name of map[string]interface{}
type M map[string]interface{}

// func (m M) String() string  {
// 	return fmt.Sprint(m)
// }

// StringMap string map short name
type StringMap map[string]string

// Level type
type Level uint32

// String get level name
func (l Level) String() string {
	return LevelName(l)
}

// Name get level name
func (l Level) Name() string {
	return LevelName(l)
}

// LowerName get lower level name
func (l Level) LowerName() string {
	return strings.ToLower(LevelName(l))
}

// Levels level list
type Levels []Level

// Contains given level
func (ls Levels) Contains(level Level) bool {
	for _, l := range ls {
		if l == level {
			return true
		}
	}
	return false
}

// BaseLoggerFace interface
type BaseLoggerFace interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

// LoggerFace interface
type LoggerFace interface {
	BaseLoggerFace
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// FullLoggerFace interface
type FullLoggerFace interface {
	LoggerFace
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
}

// FlushCloseWriter is the interface satisfied by logging destinations.
type FlushCloseWriter interface {
	Flush() error
	// the output writer
	io.WriteCloser
}

// FormatterWriterHandler interface
type FormatterWriterHandler interface {
	Handler
	// record formatter
	Formatter() Formatter
	// the output writer
	Writer() io.Writer
}

//
// Handler interface
//

// Handler interface definition
type Handler interface {
	// Close handler.
	// You should first call Flush() on close logic.
	// Refer the FileHandler.Close() handle
	io.Closer
	// Flush logs to disk
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// All records may be passed to this method, and the handler should discard
	// Those that it does not want to handle.
	Handle(*Record) error
}

//
// Processor interface
//

// Processor interface definition
type Processor interface {
	// Process record
	Process(record *Record)
}

// ProcessorFunc wrapper definition
type ProcessorFunc func(record *Record)

// Process record
func (fn ProcessorFunc) Process(record *Record) {
	fn(record)
}

// ProcessableHandler interface
type ProcessableHandler interface {
	// AddProcessor add an processor
	AddProcessor(Processor)
	// SetProcessor set the log processor
	ProcessRecord(record *Record)
}

// Processable definition
type Processable struct {
	processors []Processor
}

// AddProcessor to the handler
func (p *Processable) AddProcessor(processor Processor) {
	p.processors = append(p.processors, processor)
}

// ProcessRecord process records
func (p *Processable) ProcessRecord(r *Record) {
	// processing log record
	for _, processor := range p.processors {
		processor.Process(r)
	}
}

//
// Formatter interface
//

// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}

// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format an record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}

// FormattableHandler interface
type FormattableHandler interface {
	// Formatter get the log formatter
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}

// Formattable definition
type Formattable struct {
	formatter Formatter
}

// Formatter get formatter. if not set, will return TextFormatter
func (f *Formattable) Formatter() Formatter {
	if f.formatter == nil {
		f.formatter = NewTextFormatter()
	}
	return f.formatter
}

// SetFormatter to handler
func (f *Formattable) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}

// FormatRecord to bytes
func (f *Formattable) FormatRecord(record *Record) ([]byte, error) {
	return f.Formatter().Format(record)
}

// These are the different logging levels. You can set the logging level to log handler
const (
	// PanicLevel level, highest level of severity. will call panic() if the logging level <= PanicLevel.
	PanicLevel Level = 100
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level <= FatalLevel.
	FatalLevel Level = 200
	// ErrorLevel level. Runtime errors. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel Level = 300
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel Level = 400
	// NoticeLevel level Uncommon events
	NoticeLevel Level = 500
	// InfoLevel level. Examples: User logs in, SQL logs.
	InfoLevel Level = 600
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel Level = 700
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel Level = 800
)

var (
	FieldKeyTime = "time"
	// FieldKeyDate  = "date"
	FieldKeyData = "data"

	// NOTICE: you must set `Logger.ReportCaller=true` for reporting caller
	// FieldKeyCaller filename with line with func name. eg: "logger_test.go:48->TestLogger_ReportCaller"
	FieldKeyCaller = "caller"
	// FieldKeyFLine filename with line. eg: "logger_test.go:48"
	FieldKeyFLine = "fline"
	// FieldKeyPkg package name. "github.com/gookit/slog_test"
	FieldKeyPkg = "package"
	// FieldKeyFunc package with func name. eg: "github.com/gookit/slog_test.TestLogger_ReportCaller"
	FieldKeyFunc = "func"
	// FieldKeyFunc only func name. eg: "TestLogger_ReportCaller"
	FieldKeyFcName = "fcname"
	// FieldKeyFile full filepath with line. eg: "/work/go/gookit/slog/logger_test.go:48"
	FieldKeyFile = "file"

	FieldKeyDatetime  = "datetime"
	FieldKeyTimestamp = "timestamp"

	FieldKeyLevel = "level"
	FieldKeyError = "error"
	FieldKeyExtra = "extra"

	FieldKeyChannel = "channel"
	FieldKeyMessage = "message"
)

var (
	DefaultChannelName = "application"
	DefaultTimeFormat  = "2006/01/02 15:04:05"
	// TimeFormatRFC3339  = time.RFC3339
)

var (
	// AllLevels exposing all logging levels
	AllLevels = Levels{
		PanicLevel,
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		NoticeLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel,
	}

	DangerLevels = Levels{PanicLevel, FatalLevel, ErrorLevel, WarnLevel}
	NormalLevels = Levels{InfoLevel, NoticeLevel, DebugLevel, TraceLevel}

	// PrintLevel for use logger.Print / Printf / Println
	PrintLevel = InfoLevel

	// LevelNames all level mapping name
	LevelNames = map[Level]string{
		PanicLevel:  "PANIC",
		FatalLevel:  "FATAL",
		ErrorLevel:  "ERROR",
		NoticeLevel: "NOTICE",
		WarnLevel:   "WARNING",
		InfoLevel:   "INFO",
		DebugLevel:  "DEBUG",
		TraceLevel:  "TRACE",
	}
)

var (
	// DefaultFields default log export fields
	DefaultFields = []string{
		FieldKeyDatetime,
		FieldKeyChannel,
		FieldKeyLevel,
		FieldKeyCaller,
		FieldKeyMessage,
		FieldKeyData,
		FieldKeyExtra,
	}
	// NoTimeFields log export fields without time
	NoTimeFields = []string{
		FieldKeyChannel,
		FieldKeyLevel,
		FieldKeyMessage,
		FieldKeyData,
		FieldKeyExtra,
	}
)

// LevelName match
func LevelName(l Level) string {
	if n, ok := LevelNames[l]; ok {
		return n
	}

	return "UNKNOWN"
}

// Name2Level convert name to level
func Name2Level(ln string) (Level, error) {
	switch strings.ToLower(ln) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "err", "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "notice":
		return NoticeLevel, nil
	case "info", "": // make the zero value useful
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("invalid log Level: %q", ln)
}

// MustLevelByName convert name to level
func MustLevelByName(ln string) Level {
	l, err := Name2Level(ln)
	if err != nil {
		return InfoLevel
	}

	return l
}
