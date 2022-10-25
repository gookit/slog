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
	"time"
)

// var bufferPool *sync.Pool

func init() {
	// bufferPool = &sync.Pool{
	// 	New: func() interface{} {
	// 		return new(bytes.Buffer)
	// 	},
	// }
}

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

// Reset the std logger
func Reset() {
	ResetExitHandlers(true)
	// new std
	std = NewStdLogger()
}

// Configure the std logger
func Configure(fn func(l *SugaredLogger)) { std.Configure(fn) }

// Exit runs all the logger exit handlers and then terminates the program using os.Exit(code)
func Exit(code int) { std.Exit(code) }

// SetExitFunc to the std logger
func SetExitFunc(fn func(code int)) { std.ExitFunc = fn }

// Flush log messages
func Flush() error { return std.Flush() }

// MustFlush log messages
func MustFlush() {
	err := Flush()
	if err != nil {
		panic(err)
	}
}

// FlushTimeout flush logs with timeout.
func FlushTimeout(timeout time.Duration) { std.FlushTimeout(timeout) }

// FlushDaemon run flush handle on daemon
//
// Usage:
//
//	go slog.FlushDaemon()
func FlushDaemon() { std.FlushDaemon() }

// SetLogLevel for the std logger
func SetLogLevel(l Level) { std.Level = l }

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

// WithData new record with data
func WithData(data M) *Record {
	return std.WithData(data)
}

// WithFields new record with fields
func WithFields(fields M) *Record {
	return std.WithFields(fields)
}

// -------------------------- Add log messages with level -----------------------------

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
func Info(args ...any) {
	std.log(InfoLevel, args)
}

// Infof logs a message at level Info
func Infof(format string, args ...any) {
	std.logf(InfoLevel, format, args)
}

// Notice logs a message at level Notice
func Notice(args ...any) {
	std.log(NoticeLevel, args)
}

// Noticef logs a message at level Notice
func Noticef(format string, args ...any) {
	std.logf(NoticeLevel, format, args)
}

// Warn logs a message at level Warn
func Warn(args ...any) {
	std.log(WarnLevel, args)
}

// Warnf logs a message at level Warn
func Warnf(format string, args ...any) {
	std.logf(WarnLevel, format, args)
}

// Error logs a message at level Error
func Error(args ...any) {
	std.log(ErrorLevel, args)
}

// ErrorT logs a error type at level Error
func ErrorT(err error) {
	if err != nil {
		std.log(ErrorLevel, []any{err})
	}
}

// Errorf logs a message at level Error
func Errorf(format string, args ...any) {
	std.logf(ErrorLevel, format, args)
}

// Debug logs a message at level Debug
func Debug(args ...any) {
	std.log(DebugLevel, args)
}

// Debugf logs a message at level Debug
func Debugf(format string, args ...any) {
	std.logf(DebugLevel, format, args)
}

// Fatal logs a message at level Fatal
func Fatal(args ...any) {
	std.log(FatalLevel, args)
}

// Fatalf logs a message at level Fatal
func Fatalf(format string, args ...any) {
	std.logf(FatalLevel, format, args)
}

// FatalErr logs a message at level Fatal on err is not nil
func FatalErr(err error) {
	if err != nil {
		std.log(FatalLevel, []any{err})
	}
}

// Panic logs a message at level Panic
func Panic(args ...any) {
	std.log(PanicLevel, args)
}

// Panicf logs a message at level Panic
func Panicf(format string, args ...any) {
	std.logf(PanicLevel, format, args)
}

// PanicErr logs a message at level Panic on err is not nil
func PanicErr(err error) {
	if err != nil {
		std.log(PanicLevel, []any{err})
	}
}
