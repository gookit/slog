package slog

import (
	"fmt"
	"io"
	"strings"

)

// SugaredLogger definition
type SugaredLogger struct {
	*Logger
	out io.Writer
	Level Level
	formatter Formatter
}

// NewSugaredLogger create new SugaredLogger
func NewSugaredLogger(out io.Writer, level Level) *SugaredLogger {
	return &SugaredLogger{
		out: out,
		Level: level,
		Logger: New(),
	}
}

// SetOutput for the logger
func (sl *SugaredLogger) SetOutput(out io.Writer) {
	sl.out = out
}

// Formatter get formatter
func (sl *SugaredLogger) Formatter() slog.Formatter {
	if sl.formatter == nil {
		sl.formatter = DefaultFormatter
	}
	return sl.formatter
}

// SetFormatter to handler
func (sl *SugaredLogger) SetFormatter(formatter slog.Formatter) {
	sl.formatter = formatter
}

// IsHandling Check if the current level can be handling
func (sl *SugaredLogger) IsHandling(level Level) bool {
	return sl.Level >= level
}

// Handle log record
func (sl *SugaredLogger) Handle(record *slog.Record)  error {
	bts, err := sl.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = h.Out.Write(bts)
	return err
}

var std = NewWithName("stdLogger")

// Std get std logger
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

// Exit runs all the logger exit handlers and then terminates the program using os.Exit(code)
func Exit(code int) {
	std.Exit(code)
}

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