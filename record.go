package slog

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gookit/goutil/strutil"
)

// Record a log record definition
type Record struct {
	logger *Logger

	Time  time.Time
	Level Level
	// level name cache from Level
	levelName string

	// Channel log channel name. eg: "order", "goods", "user"
	Channel string
	Message string

	// Ctx context.Context
	Ctx context.Context

	// Buffer Can use Buffer on formatter
	Buffer *bytes.Buffer

	// Fields custom fields data.
	// Contains all the fields set by the user.
	Fields M

	// Data log context data
	Data M

	// Extra log extra data
	Extra M

	// Caller information
	Caller *runtime.Frame
	// Formatted []byte

	// stacks []byte
	// cache the r.Time.Nanosecond() / 1000
	microSecond int
	// field caches mapping for optimize performance. TODO use map[string][]byte ?
	strmp map[string]string
}

var (
	// ErrorKey Define the key when adding errors using WithError.
	ErrorKey   = "error"
	bufferPool *sync.Pool
)

func newRecord(logger *Logger) *Record {
	return &Record{
		logger:  logger,
		Channel: DefaultChannelName,
		// init map data field
		// Data:   make(M, 2),
		// Extra:  make(M, 0),
		// Fields: make(M, 0),
	}
}

//
// ---------------------------------------------------------------------------
// Copy record with something
// ---------------------------------------------------------------------------
//

// WithTime set the record time
func (r *Record) WithTime(t time.Time) *Record {
	nr := r.Copy()
	nr.Time = t
	return nr
}

// WithContext on record
func (r *Record) WithContext(ctx context.Context) *Record {
	nr := r.Copy()
	nr.Ctx = ctx
	return nr
}

// WithError on record
func (r *Record) WithError(err error) *Record {
	return r.WithFields(M{ErrorKey: err})
}

// WithData on record
func (r *Record) WithData(data M) *Record {
	nr := r.Copy()
	// if nr.Data == nil {
	// 	nr.Data = data
	// }

	nr.Data = data
	return nr
}

// WithField with an new field to record
func (r *Record) WithField(name string, val interface{}) *Record {
	return r.WithFields(M{name: val})
}

// WithFields with new fields to record
func (r *Record) WithFields(fields M) *Record {
	nr := r.Copy()
	if nr.Fields == nil {
		nr.Fields = make(M, len(fields))
	}

	for k, v := range fields {
		nr.Fields[k] = v
	}
	return nr
}

// Copy new record from old record
func (r *Record) Copy() *Record {
	dataCopy := make(M, len(r.Data))
	for k, v := range r.Data {
		dataCopy[k] = v
	}

	fieldsCopy := make(M, len(r.Fields))
	for k, v := range r.Fields {
		fieldsCopy[k] = v
	}

	extraCopy := make(M, len(r.Extra))
	for k, v := range r.Extra {
		extraCopy[k] = v
	}

	return &Record{
		logger:    r.logger,
		Channel:   r.Channel,
		Time:      r.Time,
		Level:     r.Level,
		levelName: r.levelName,
		Message:   r.Message,
		Data:      dataCopy,
		Extra:     extraCopy,
		Fields:    fieldsCopy,
	}
}

//
// ---------------------------------------------------------------------------
// Direct set value to record
// ---------------------------------------------------------------------------
//

// SetContext on record
func (r *Record) SetContext(ctx context.Context) *Record {
	r.Ctx = ctx
	return r
}

// SetData on record
func (r *Record) SetData(data M) *Record {
	r.Data = data
	return r
}

// AddData on record
func (r *Record) AddData(data M) *Record {
	if r.Data == nil {
		r.Data = data
		return r
	}

	for k, v := range data {
		r.Data[k] = v
	}
	return r
}

// AddValue add Data value to record
func (r *Record) AddValue(key string, value interface{}) *Record {
	if r.Data == nil {
		r.Data = make(M, 8)
		return r
	}

	r.Data[key] = value
	return r
}

// SetExtra information on record
func (r *Record) SetExtra(data M) *Record {
	r.Extra = data
	return r
}

// AddExtra information on record
func (r *Record) AddExtra(data M) *Record {
	if r.Extra == nil {
		r.Extra = data
		return r
	}

	for k, v := range data {
		r.Extra[k] = v
	}
	return r
}

// SetExtraValue on record
func (r *Record) SetExtraValue(k string, v interface{}) {
	if r.Extra == nil {
		r.Extra = make(M, 8)
	}

	r.Extra[k] = v
}

// SetTime on record
func (r *Record) SetTime(t time.Time) *Record {
	r.Time = t
	return r
}

// AddField add new field to the record
func (r *Record) AddField(name string, val interface{}) *Record {
	if r.Fields == nil {
		r.Fields = make(M, 8)
	}

	r.Fields[name] = val
	return r
}

// AddFields add new fields to the record
func (r *Record) AddFields(fields M) *Record {
	if r.Fields == nil {
		r.Fields = fields
		return r
	}

	for n, v := range fields {
		r.Fields[n] = v
	}
	return r
}

// SetFields to the record
func (r *Record) SetFields(fields M) *Record {
	r.Fields = fields
	return r
}

// Object data on record TODO optimize performance
// func (r *Record) Object(obj fmt.Stringer) *Record {
// 	r.Data = ctx
// 	return r
// }

//
// ---------------------------------------------------------------------------
// Add log message with level
// ---------------------------------------------------------------------------
//

// Log an message with level
func (r *Record) Log(level Level, args ...interface{}) {
	r.logBytes(level, formatArgsWithSpaces(args))
}

// Logf an message with level
func (r *Record) Logf(level Level, format string, args ...interface{}) {
	r.logBytes(level, []byte(fmt.Sprintf(format, args...)))
}

// Info logs a message at level Info
func (r *Record) Info(args ...interface{}) {
	r.Log(InfoLevel, args...)
}

// Infof logs a message at level Info
func (r *Record) Infof(format string, args ...interface{}) {
	r.Logf(InfoLevel, format, args...)
}

// Trace logs a message at level Trace
func (r *Record) Trace(args ...interface{}) {
	r.Log(TraceLevel, args...)
}

// Tracef logs a message at level Trace
func (r *Record) Tracef(format string, args ...interface{}) {
	r.Logf(TraceLevel, format, args...)
}

// Error logs a message at level Error
func (r *Record) Error(args ...interface{}) {
	r.Log(ErrorLevel, args...)
}

// Errorf logs a message at level Error
func (r *Record) Errorf(format string, args ...interface{}) {
	r.Logf(ErrorLevel, format, args...)
}

// Notice logs a message at level Notice
func (r *Record) Notice(args ...interface{}) {
	r.Log(NoticeLevel, args...)
}

// Noticef logs a message at level Notice
func (r *Record) Noticef(format string, args ...interface{}) {
	r.Logf(NoticeLevel, format, args...)
}

// Debug logs a message at level Debug
func (r *Record) Debug(args ...interface{}) {
	r.Log(DebugLevel, args...)
}

// Debugf logs a message at level Debug
func (r *Record) Debugf(format string, args ...interface{}) {
	r.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func (r *Record) Fatal(args ...interface{}) {
	r.Log(FatalLevel, args...)
}

// Fatalf logs a message at level Fatal
func (r *Record) Fatalf(format string, args ...interface{}) {
	r.Logf(FatalLevel, format, args...)
}

// Panic logs a message at level Panic
func (r *Record) Panic(args ...interface{}) {
	r.Log(PanicLevel, args...)
}

// Panicf logs a message at level Panic
func (r *Record) Panicf(format string, args ...interface{}) {
	r.Logf(PanicLevel, format, args...)
}

// ---------------------------------------------------------------------------
// helper methods
// ---------------------------------------------------------------------------

// NewBuffer get or create an Buffer
func (r *Record) NewBuffer() *bytes.Buffer {
	if r.Buffer == nil {
		return &bytes.Buffer{}
	}

	return r.Buffer
}

// LevelName get
func (r *Record) LevelName() string {
	return r.levelName
}

// MicroSecond of the record
func (r *Record) MicroSecond() int {
	return r.microSecond
}

// func (r *Record) logString(level Level, message string) {
// 	r.logBytes(level, []byte(message))
// }

func (r *Record) logBytes(level Level, message []byte) {
	r.Level = level
	// r.Message = string(message)
	// Will reduce memory allocation once
	r.Message = strutil.Byte2str(message)

	// var buffer *bytes.Buffer
	// buffer = bufferPool.Get().(*bytes.Buffer)
	// buffer.Reset()
	// defer bufferPool.Put(buffer)

	// r.Buffer = buffer
	r.initLogTime()

	// do write log message
	r.logger.write(level, r)

	// r.Buffer = nil
	// TODO release on here ??
	// r.logger.releaseRecord(r)
}

func (r *Record) initLogTime() {
	if r.Time.IsZero() {
		r.Time = time.Now()
	}

	r.microSecond = r.Time.Nanosecond() / 1000
}
