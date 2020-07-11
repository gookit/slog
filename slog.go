package slog

import (
	"fmt"
	"os"
	"strings"
)

// Level type
type Level uint32

// String get level name
func (l Level) String() string  {
	return LevelName(l)
}

// Name get level name
func (l Level) Name() string  {
	return LevelName(l)
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Runtime errors. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// NoticeLevel level Uncommon events
	NoticeLevel
	// InfoLevel level. Examples: User logs in, SQL logs.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// A constant exposing all logging levels
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	NoticeLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

var LevelNames = map[Level]string{
	PanicLevel:  "PANIC",
	FatalLevel:  "FATAL",
	ErrorLevel:  "ERROR",
	NoticeLevel: "NOTICE",
	WarnLevel:   "WARNING",
	InfoLevel:   "INFO",
	DebugLevel:  "DEBUG",
	TraceLevel:  "TRACE",
}

var std = New("stdLogger")

var (

)

func Std() *Logger  {
	return std
}

func AddHandler(h Handler) {
	std.AddHandler(h)
}

func AddProcessor(p Processor) {
	std.AddProcessor(p)
}

// Trace logs a message at level Trace
func Trace(args ...interface{}) {
	std.Log(TraceLevel, args...)
}

// Trace logs a message at level Trace
func Tracef(format string, args ...interface{})  {
	std.Logf(TraceLevel, format, args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	std.Log(InfoLevel, args...)
}

// Info logs a message at level Info
func Infof(format string, args ...interface{})  {
	std.Logf(InfoLevel, format, args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	std.Log(WarnLevel, args...)
}

// Warn logs a message at level Warn
func Warnf(format string, args ...interface{})  {
	std.Logf(WarnLevel, format, args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	std.Log(ErrorLevel, args...)
}

// Error logs a message at level Error
func Errorf(format string, args ...interface{})  {
	std.Logf(ErrorLevel, format, args...)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	std.Log(DebugLevel, args...)
}

// Debug logs a message at level Debug
func Debugf(format string, args ...interface{})  {
	std.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func Fatal(args ...interface{}) {
	std.Log(FatalLevel, args...)
}

// Fatal logs a message at level Fatal
func Fatalf(format string, args ...interface{})  {
	std.Logf(FatalLevel, format, args...)
}

// Exit runs all the Logrus atexit handlers and then terminates the program using os.Exit(code)
func Exit(code int) {
	std.runExitHandlers()
	runExitHandlers()
	os.Exit(code)
}

func LevelName(l Level) string {
	if n, ok := LevelNames[l]; ok {
		return n
	}

	return "UNKNOWN"
}

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
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("invalid log Level: %q", ln)
}