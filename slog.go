/*
Package slog Lightweight, extensible, configurable logging library written in Go.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/slog

Quick usage:

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

More usage please see README.
*/
package slog

import (
	"context"
	"time"

	"github.com/gookit/goutil"
	"github.com/gookit/gsr"
)

// SLogger interface
type SLogger interface {
	gsr.Logger
	Log(level Level, v ...any)
	Logf(level Level, format string, v ...any)
}

// LoggerFn func
type LoggerFn func(l *Logger)

//
// ------------------------------------------------------------
// Global std logger operate
// ------------------------------------------------------------
//

// std logger is an SugaredLogger.
// It is directly available without any additional configuration
var std = NewStdLogger()

// Std get std logger
func Std() *SugaredLogger { return std }

// Reset the std logger and reset exit handlers
func Reset() {
	ResetExitHandlers(true)
	// new std
	std = NewStdLogger()
}

// Configure the std logger
func Configure(fn func(l *SugaredLogger)) { std.Config(fn) }

// SetExitFunc to the std logger
func SetExitFunc(fn func(code int)) { std.ExitFunc = fn }

// Exit runs all exit handlers and then terminates the program using os.Exit(code)
func Exit(code int) { std.Exit(code) }

// Close logger, flush and close all handlers.
//
// IMPORTANT: please call Close() before app exit.
func Close() error { return std.Close() }

// MustClose logger, flush and close all handlers.
//
// IMPORTANT: please call Close() before app exit.
func MustClose() { goutil.PanicErr(Close()) }

// Flush log messages
func Flush() error { return std.Flush() }

// MustFlush log messages
func MustFlush() { goutil.PanicErr(Flush()) }

// FlushTimeout flush logs with timeout.
func FlushTimeout(timeout time.Duration) { std.FlushTimeout(timeout) }

// FlushDaemon run flush handle on daemon.
//
// Usage please see slog_test.ExampleFlushDaemon()
func FlushDaemon(onStops ...func()) {
	std.FlushDaemon(onStops...)
}

// StopDaemon stop flush daemon
func StopDaemon() { std.StopDaemon() }

// SetLogLevel max level for the std logger
func SetLogLevel(l Level) { std.Level = l }

// SetLevelByName set max log level by name. eg: "info", "debug" ...
func SetLevelByName(name string) { std.Level = LevelByName(name) }

// SetFormatter to std logger
func SetFormatter(f Formatter) { std.Formatter = f }

// GetFormatter of the std logger
func GetFormatter() Formatter { return std.Formatter }

// AddHandler to the std logger
func AddHandler(h Handler) { std.AddHandler(h) }

// PushHandler to the std logger
func PushHandler(h Handler) { std.AddHandler(h) }

// AddHandlers to the std logger
func AddHandlers(hs ...Handler) { std.AddHandlers(hs...) }

// PushHandlers to the std logger
func PushHandlers(hs ...Handler) { std.PushHandlers(hs...) }

// AddProcessor to the logger
func AddProcessor(p Processor) { std.AddProcessor(p) }

// AddProcessors to the logger
func AddProcessors(ps ...Processor) { std.AddProcessors(ps...) }

// -------------------------- New record with log data, fields -----------------------------

// WithExtra new record with extra data
func WithExtra(ext M) *Record {
	return std.WithExtra(ext)
}

// WithData new record with data
func WithData(data M) *Record {
	return std.WithData(data)
}

// WithValue new record with data value
func WithValue(key string, value any) *Record {
	return std.WithValue(key, value)
}

// WithField new record with field.
//
// TIP: add field need config Formatter template fields.
func WithField(name string, value any) *Record {
	return std.WithField(name, value)
}

// WithFields new record with fields
//
// TIP: add field need config Formatter template fields.
func WithFields(fields M) *Record {
	return std.WithFields(fields)
}

// WithContext new record with context
func WithContext(ctx context.Context) *Record {
	return std.WithContext(ctx)
}

// -------------------------- Add log messages with level -----------------------------

// Log logs a message with level
func Log(level Level, args ...any) { std.log(level, args) }

// Print logs a message at level PrintLevel
func Print(args ...any) { std.log(PrintLevel, args) }

// Println logs a message at level PrintLevel
func Println(args ...any) { std.log(PrintLevel, args) }

// Printf logs a message at level PrintLevel
func Printf(format string, args ...any) { std.logf(PrintLevel, format, args) }

// Trace logs a message at level Trace
func Trace(args ...any) { std.log(TraceLevel, args) }

// Tracef logs a message at level Trace
func Tracef(format string, args ...any) { std.logf(TraceLevel, format, args) }

// Info logs a message at level Info
func Info(args ...any) { std.log(InfoLevel, args) }

// Infof logs a message at level Info
func Infof(format string, args ...any) { std.logf(InfoLevel, format, args) }

// Notice logs a message at level Notice
func Notice(args ...any) { std.log(NoticeLevel, args) }

// Noticef logs a message at level Notice
func Noticef(format string, args ...any) { std.logf(NoticeLevel, format, args) }

// Warn logs a message at level Warn
func Warn(args ...any) { std.log(WarnLevel, args) }

// Warnf logs a message at level Warn
func Warnf(format string, args ...any) { std.logf(WarnLevel, format, args) }

// Error logs a message at level Error
func Error(args ...any) { std.log(ErrorLevel, args) }

// Errorf logs a message at level Error
func Errorf(format string, args ...any) { std.logf(ErrorLevel, format, args) }

// ErrorT logs a error type at level Error
func ErrorT(err error) {
	if err != nil {
		std.log(ErrorLevel, []any{err})
	}
}

// EStack logs a error message and with call stack.
// func EStack(args ...any) {
// 	std.WithExtra(map[string]any{"stack": goinfo.GetCallerInfo(2)}).
// 		log(ErrorLevel, args)
// }

// Debug logs a message at level Debug
func Debug(args ...any) { std.log(DebugLevel, args) }

// Debugf logs a message at level Debug
func Debugf(format string, args ...any) { std.logf(DebugLevel, format, args) }

// Fatal logs a message at level Fatal
func Fatal(args ...any) { std.log(FatalLevel, args) }

// Fatalf logs a message at level Fatal
func Fatalf(format string, args ...any) { std.logf(FatalLevel, format, args) }

// FatalErr logs a message at level Fatal on err is not nil
func FatalErr(err error) {
	if err != nil {
		std.log(FatalLevel, []any{err})
	}
}

// Panic logs a message at level Panic
func Panic(args ...any) { std.log(PanicLevel, args) }

// Panicf logs a message at level Panic
func Panicf(format string, args ...any) { std.logf(PanicLevel, format, args) }

// PanicErr logs a message at level Panic on err is not nil
func PanicErr(err error) {
	if err != nil {
		std.log(PanicLevel, []any{err})
	}
}
