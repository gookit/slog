package slog

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
	PanicLevel Level = iota + 1
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

var (
	DefaultChannelName = "application"
	DefaultTimeFormat  = "2006/01/02 15:04:05"
)

const (
	FieldKeyTime  = "time"
	FieldKeyData  = "data"
	FieldKeyFunc  = "func"
	FieldKeyFile  = "file"
	// FieldKeyDate  = "date"

	FieldKeyDatetime  = "datetime"

	FieldKeyLevel = "level"
	FieldKeyError = "error"
	FieldKeyExtra = "extra"

	FieldKeyChannel = "channel"
	FieldKeyMessage = "message"
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
