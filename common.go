package slog

import (
	"errors"
	"strings"
	"time"
)

//
// log level definitions
//

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
	if n, ok := lowerLevelNames[l]; ok {
		return n
	}
	return "unknown"
}

// ShouldHandling compare level, if current level <= l, it will be record.
func (l Level) ShouldHandling(curLevel Level) bool {
	return curLevel <= l
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

// These are the different logging levels. You can set the logging level to log handler
const (
	// PanicLevel level, the highest level of severity. will call panic() if the logging level <= PanicLevel.
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

//
// some commonly definitions
//

// StringMap string map short name
type StringMap = map[string]string

// M short name of map[string]any
type M map[string]any

// String map to string
func (m M) String() string {
	return mapToString(m)
}

// ClockFn func
type ClockFn func() time.Time

// Now implements the Clocker
func (fn ClockFn) Now() time.Time {
	return fn()
}

// NOTICE: you must set `Logger.ReportCaller=true` for reporting caller.
// then config the Logger.CallerFlag by follow flags.
const (
	// CallerFlagFnlFcn report short func name with filename and with line.
	// eg: "logger_test.go:48,TestLogger_ReportCaller"
	CallerFlagFnlFcn uint8 = iota
	// CallerFlagFull full func name with filename and with line.
	// eg: "github.com/gookit/slog_test.TestLogger_ReportCaller(),logger_test.go:48"
	CallerFlagFull
	// CallerFlagFunc full package with func name.
	// eg: "github.com/gookit/slog_test.TestLogger_ReportCaller"
	CallerFlagFunc
	// CallerFlagFcLine full package with func name and with line.
	// eg: "github.com/gookit/slog_test.TestLogger_ReportCaller:48"
	CallerFlagFcLine
	// CallerFlagPkg report full package name.
	// eg: "github.com/gookit/slog_test"
	CallerFlagPkg
	// CallerFlagPkgFnl report full package name + filename + line.
	// eg: "github.com/gookit/slog_test,logger_test.go:48"
	CallerFlagPkgFnl
	// CallerFlagFpLine report full filepath with line.
	// eg: "/work/go/gookit/slog/logger_test.go:48"
	CallerFlagFpLine
	// CallerFlagFnLine report filename with line.
	// eg: "logger_test.go:48"
	CallerFlagFnLine
	// CallerFlagFcName only report func name.
	// eg: "TestLogger_ReportCaller"
	CallerFlagFcName
)

var (
	// FieldKeyData define the key name for Record.Data
	FieldKeyData = "data"
	// FieldKeyTime key name
	FieldKeyTime = "time"
	// FieldKeyDate key name
	FieldKeyDate = "date"

	// FieldKeyDatetime key name
	FieldKeyDatetime = "datetime"
	// FieldKeyTimestamp key name
	FieldKeyTimestamp = "timestamp"

	// FieldKeyCaller the field key name for report caller.
	//
	// For caller style please see CallerFlagFull, CallerFlagFunc and more.
	//
	// NOTICE: you must set `Logger.ReportCaller=true` for reporting caller
	FieldKeyCaller = "caller"

	// FieldKeyLevel name
	FieldKeyLevel = "level"
	// FieldKeyError Define the key when adding errors using WithError.
	FieldKeyError = "error"
	// FieldKeyExtra key name
	FieldKeyExtra = "extra"

	// FieldKeyChannel name
	FieldKeyChannel = "channel"
	// FieldKeyMessage name
	FieldKeyMessage = "message"
)

var (
	// DefaultChannelName for log record
	DefaultChannelName = "application"
	// DefaultTimeFormat define
	DefaultTimeFormat = "2006/01/02T15:04:05.000"
	// TimeFormatRFC3339  = time.RFC3339

	// DoNothingOnExit handle func. use for testing.
	DoNothingOnExit = func(code int) {}
	// DoNothingOnPanic handle func. use for testing.
	DoNothingOnPanic = func(v any) {}

	// DefaultPanicFn handle func
	DefaultPanicFn = func(v any) {
		panic(v)
	}
	// DefaultClockFn create func
	DefaultClockFn = ClockFn(func() time.Time {
		return time.Now()
	})
)

var (
	// PrintLevel for use logger.Print / Printf / Println
	PrintLevel = InfoLevel

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

	// DangerLevels define the commonly danger log levels
	DangerLevels = Levels{PanicLevel, FatalLevel, ErrorLevel, WarnLevel}
	// NormalLevels define the commonly normal log levels
	NormalLevels = Levels{InfoLevel, NoticeLevel, DebugLevel, TraceLevel}

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

	// lower level name.
	lowerLevelNames = buildLowerLevelName()
	// empty time for reset record.
	emptyTime = time.Time{}
)

// LevelName match
func LevelName(l Level) string {
	if n, ok := LevelNames[l]; ok {
		return n
	}
	return "UNKNOWN"
}

// LevelByName convert name to level
func LevelByName(ln string) Level {
	l, err := Name2Level(ln)
	if err != nil {
		return InfoLevel
	}
	return l
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
	return 0, errors.New("invalid log level name: " + ln)
}

//
// exit handle logic
//

// global exit handler
var exitHandlers = make([]func(), 0)

func runExitHandlers() {
	defer func() {
		if err := recover(); err != nil {
			printlnStderr("slog: run exit handler error:", err)
		}
	}()

	for _, handler := range exitHandlers {
		handler()
	}
}

// ExitHandlers get all global exitHandlers
func ExitHandlers() []func() {
	return exitHandlers
}

// RegisterExitHandler register an exit-handler on global exitHandlers
func RegisterExitHandler(handler func()) {
	exitHandlers = append(exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func PrependExitHandler(handler func()) {
	exitHandlers = append([]func(){handler}, exitHandlers...)
}

// ResetExitHandlers reset all exitHandlers
func ResetExitHandlers(applyToStd bool) {
	exitHandlers = make([]func(), 0)

	if applyToStd {
		std.ResetExitHandlers()
	}
}

func (l *Logger) runExitHandlers() {
	defer func() {
		if err := recover(); err != nil {
			printlnStderr("slog: run exit handler error:", err)
		}
	}()

	for _, handler := range l.exitHandlers {
		handler()
	}
}

// RegisterExitHandler register an exit-handler on global exitHandlers
func (l *Logger) RegisterExitHandler(handler func()) {
	l.exitHandlers = append(l.exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func (l *Logger) PrependExitHandler(handler func()) {
	l.exitHandlers = append([]func(){handler}, l.exitHandlers...)
}

// ResetExitHandlers reset logger exitHandlers
func (l *Logger) ResetExitHandlers() {
	l.exitHandlers = make([]func(), 0)
}

// ExitHandlers get all exitHandlers of the logger
func (l *Logger) ExitHandlers() []func() {
	return l.exitHandlers
}
