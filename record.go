package slog

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// StringMap string map
type StringMap map[string]string

// Record a log record definition
type Record struct {
	logger *Logger

	Time  time.Time
	Level Level

	LevelName string

	// Channel log channel name. eg: "order", "goods", "user"
	Channel string
	Message string

	// Ctx context.Context
	Ctx context.Context

	// Fields custom fields. Contains all the fields set by the user.
	Fields M

	// log data
	Data M

	// Extra log extra data
	Extra M

	// Caller information
	Caller *runtime.Frame
	// Formatted []byte
}

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

var (
	// Defines the key when adding errors using WithError.
	ErrorKey = "error"

	bufferPool *sync.Pool
)

// DefaultFields default log export fields
var DefaultFields = []string{
	FieldKeyDatetime,
	FieldKeyChannel,
	FieldKeyLevel,
	FieldKeyMessage,
	FieldKeyData,
	FieldKeyExtra,
}

func newRecord(logger *Logger) *Record {
	return &Record{
		logger:  logger,
		Fields:  make(M),
		Channel: "application",
	}
}

// WithContext on record
func (r *Record) WithContext(ctx context.Context) *Record {
	r.Ctx = ctx

	return r
}

func (r *Record) WithError(err error) *Record {
	return r.WithFields(M{ErrorKey: err})
}

// WithField with an new field to record
func (r *Record) WithField(name string, val interface{}) *Record {
	return r.WithFields(M{name: val})
}

// WithField with new fields to record
func (r *Record) WithFields(fields M) *Record {
	// data := make(M, len(r.Data)+len(fields))
	// for k, v := range r.Data {
	// 	data[k] = v
	// }

	return &Record{
		logger:    r.logger,
		Channel:   r.Channel,
		Time:      r.Time,
		Level:     r.Level,
		LevelName: r.LevelName,
		Message:   r.Message,
		Fields:    fields,
	}
}

// AddField add new field to the record
func (r *Record) AddField(name string, val interface{}) *Record {
	r.Fields[name] = val
	return r
}

// AddFields add new fields to the record
func (r *Record) AddFields(fields M) *Record {
	for n, v := range fields {
		r.Fields[n] = v
	}
	return r
}

// Log an message with level
func (r *Record) Log(level Level, args ...interface{}) {
	r.log(level, fmt.Sprint(args...))
}

// Log an message with level
func (r *Record) Logf(level Level, format string, args ...interface{}) {
	r.log(level, fmt.Sprintf(format, args...))
}

func (r *Record) log(level Level, message string) {
	r.Level = level
	r.LevelName = level.String()
	r.Message = message

	// TODO
	r.logger.write(level, r)
}
