package slog

import (
	"fmt"
	"strings"
	"time"
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

// These are the different logging levels. You can set the logging level to log handler
const (
	// PanicLevel level, highest level of severity.
	PanicLevel Level = 100
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
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

const flushInterval = 30 * time.Second

const (
	FieldKeyTime = "time"
	FieldKeyData = "data"
	FieldKeyFunc = "func"
	FieldKeyFile = "file"
	// FieldKeyDate  = "date"

	FieldKeyDatetime = "datetime"

	FieldKeyLevel = "level"
	FieldKeyError = "error"
	FieldKeyExtra = "extra"

	FieldKeyChannel = "channel"
	FieldKeyMessage = "message"
)

var (
	DefaultChannelName = "application"
	DefaultTimeFormat  = "2006/01/02 15:04:05"
)

// AllLevels exposing all logging levels
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

// LevelNames all level mapping name
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

// DefaultFields default log export fields
var DefaultFields = []string{
	FieldKeyDatetime,
	FieldKeyChannel,
	FieldKeyLevel,
	FieldKeyMessage,
	FieldKeyData,
	FieldKeyExtra,
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
