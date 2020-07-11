package slog

import (
	"os"
	"sync"
	"time"
)

type M map[string]interface{}

// Logger definition
type Logger struct {
	name string

	handlers   []Handler
	processors []Processor

	// timezone
	tz time.Time

	mu sync.Mutex

	// Reusable empty record
	recordPool sync.Pool

	exitHandlers []func()
	ExitFunc     func(int)
}

// New create an new logger
func New(name string) *Logger {
	logger := &Logger{
		name:         name,
		// exit handle
		ExitFunc:     os.Exit,
		exitHandlers: []func(){},
	}

	logger.recordPool.New = func() interface{} {
		return newRecord(logger)
	}

	return logger
}

func (logger *Logger) newRecord() *Record {
	return logger.recordPool.Get().(*Record)
}

func (logger *Logger) releaseRecord(r *Record) {
	logger.recordPool.Put(r)
}

// AddProcessor to the logger
func (logger *Logger) AddProcessor(p Processor) {
	logger.processors = append(logger.processors, p)
}

// AddProcessors to the logger
func (logger *Logger) AddProcessors(ps []Processor) {
	logger.processors = append(logger.processors, ps...)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	r := logger.newRecord()
	r.Log(level, args...)

	logger.releaseRecord(r)
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	r := logger.newRecord()
	r.Logf(level, format, args...)

	logger.releaseRecord(r)
}

func (logger *Logger) WithFields(fields M) *Record {
	r := logger.newRecord()
	defer logger.releaseRecord(r)

	return r.WithFields(fields)
}

func (logger *Logger) Exit(code int) {
	logger.runExitHandlers()

	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}

// func (logger *Logger) addRecord(level Level, message string, extra M) {
//
// }

func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Name() string {
	return logger.name
}
