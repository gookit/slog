package slog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Record a log record
type Record struct {
	logger *Logger

	Level     Level
	LevelName string
	Channel   string
	Message   string

	Time  time.Time

	Ctx context.Context

	// Contains all the fields set by the user.
	Fields M

	// log data
	Data M

	// log extra data
	Extra M

	formatted []byte
}

var (
	// Defines the key when adding errors using WithError.
	ErrorKey   = "error"

	bufferPool *sync.Pool
)

func newRecord(logger *Logger) *Record {
	return &Record{logger: logger}
}

func (r *Record) WithContext(ctx context.Context) *Record {
	r.Ctx = ctx

	return r
}

func (r *Record) WithError(err error) *Record  {
	return r.WithFields(M{ErrorKey: err})
}

func (r *Record) WithField(name string, val interface{}) *Record  {
	return r.WithFields(M{name: val})
}

func (r *Record) WithFields(fields M) *Record  {
	// data := make(M, len(r.Data)+len(fields))
	// for k, v := range r.Data {
	// 	data[k] = v
	// }

	return &Record{
		logger:    r.logger,
		Time:  r.Time,
		Level:     r.Level,
		LevelName: r.LevelName,
		Message:   r.Message,
		Fields:   fields,
	}
}

// AddField add new field to the record
func (r *Record) AddField(name string, val interface{}) *Record  {
	r.Fields[name] = val
	return r
}

func (r *Record) Log(level Level, args ...interface{}){
	r.Level = level
	r.LevelName = level.String()
	r.Message = fmt.Sprint(args...)
}

func (r *Record) Logf(level Level, format string, args ...interface{}){
	r.Level = level
	r.LevelName = level.String()
	r.Message = fmt.Sprintf(format, args...)
}
