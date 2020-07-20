package slog

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger definition
type Logger struct {
	name string

	handlers   []Handler
	processors []Processor

	// timezone
	tz time.Time

	mu sync.Mutex

	ReportCaller bool
	MaxCallerDepth int

	// Reusable empty record
	recordPool sync.Pool

	exitHandlers []func()
	ExitFunc     func(code int)
}

// New create an new logger
func New() *Logger {
	return NewWithName("")
}

// NewWithConfig create an new logger with config func
func NewWithConfig(fn func(logger *Logger)) *Logger {
	return New().Configure(fn)
}

// NewWithHandlers create an new logger with handlers
func NewWithHandlers(hs ...Handler) *Logger {
	logger := NewWithName("")
	logger.AddHandlers(hs...)
	return logger
}

// NewWithName create an new logger with name
func NewWithName(name string) *Logger {
	logger := &Logger{
		name: name,
		// exit handle
		ExitFunc:     os.Exit,
		exitHandlers: []func(){},
		// options
		MaxCallerDepth: defaultMaxCallerDepth,
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
	// reset data
	r.Data = map[string]interface{}{}
	r.Extra = map[string]interface{}{}
	r.Fields = map[string]interface{}{}
	logger.recordPool.Put(r)
}

//
// ---------------------------------------------------------------------------
// Management logger
// ---------------------------------------------------------------------------
//

// Configure current logger
func (logger *Logger) Configure(fn func(logger *Logger)) *Logger {
	fn(logger)
	return logger
}

// Sync flushes buffered logs (if any).
func (logger *Logger) Sync() error {
	// TODO
	return nil
}

// FlushDaemon run flush handle on daemon
//
// Usage:
// 	go slog.FlushDaemon()
func (logger *Logger) FlushDaemon() {
	for range time.NewTicker(flushInterval).C {
		logger.lockAndFlushAll()
	}
}

// lockAndFlushAll is like flushAll but locks l.mu first.
func (logger *Logger) lockAndFlushAll() {
	logger.mu.Lock()
	logger.FlushAll()
	logger.mu.Unlock()
}

// FlushAll flushes all the logs and attempts to "sync" their data to disk.
// logger.mu is held.
func (logger *Logger) FlushAll() {
	// Flush from fatal down, in case there's trouble flushing.
	for _, handler := range logger.handlers {
		_= handler.Flush() // ignore error
	}
}

// Reset the logger
func (logger *Logger) Reset() {
	logger.ResetHandlers()
	logger.ResetProcessors()
}

// ResetProcessors for the logger
func (logger *Logger) ResetProcessors() {
	logger.processors = make([]Processor, 0)
}

// ResetHandlers for the logger
func (logger *Logger) ResetHandlers() {
	logger.handlers = make([]Handler, 0)
}

// Exit logger handle
func (logger *Logger) Exit(code int) {
	logger.runExitHandlers()

	// global exit handlers
	runExitHandlers()

	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}

// SetName for logger
func (logger *Logger) SetName(name string) *Logger {
	logger.name = name
	return logger
}

// Name of the logger
func (logger *Logger) Name() string {
	return logger.name
}

//
// ---------------------------------------------------------------------------
// Register handlers and processors
// ---------------------------------------------------------------------------
//

// AddHandler to the logger
func (logger *Logger) AddHandler(h Handler) {
	logger.handlers = append(logger.handlers, h)
}

// AddHandlers to the logger
func (logger *Logger) AddHandlers(hs ...Handler) {
	logger.handlers = append(logger.handlers, hs...)
}

// PushHandler to the logger. alias of AddHandler()
func (logger *Logger) PushHandler(h Handler) {
	logger.AddHandler(h)
}

// SetHandlers for the logger
func (logger *Logger) SetHandlers(hs []Handler)  {
	logger.handlers = hs
}

// AddProcessor to the logger
func (logger *Logger) AddProcessor(p Processor) {
	logger.processors = append(logger.processors, p)
}

// PushProcessor to the logger. alias of AddProcessor()
func (logger *Logger) PushProcessor(p Processor) {
	logger.processors = append(logger.processors, p)
}

// AddProcessors to the logger
func (logger *Logger) AddProcessors(ps ...Processor) {
	logger.processors = append(logger.processors, ps...)
}

// SetProcessors for the logger
func (logger *Logger) SetProcessors(ps []Processor) {
	logger.processors = ps
}

//
// ---------------------------------------------------------------------------
// New record with log data, fields
// ---------------------------------------------------------------------------
//

// SetFields new record with fields
func (logger *Logger) WithFields(fields M) *Record {
	r := logger.newRecord()
	defer logger.releaseRecord(r)

	return r.WithFields(fields)
}

// WithData new record with data
func (logger *Logger) WithData(data M) *Record {
	r := logger.newRecord()
	defer logger.releaseRecord(r)

	return r.WithData(data)
}

// WithTime new record with time.Time
func (logger *Logger) WithTime(t time.Time) *Record {
	r := logger.newRecord()
	defer logger.releaseRecord(r)

	return r.WithTime(t)
}

// WithContext new record with context.Context
func (logger *Logger) WithContext(ctx context.Context) *Record {
	r := logger.newRecord()
	defer logger.releaseRecord(r)

	return r.WithContext(ctx)
}

//
// ---------------------------------------------------------------------------
// Add log message with level
// ---------------------------------------------------------------------------
//

// Log an message
func (logger *Logger) Log(level Level, args ...interface{}) {
	r := logger.newRecord()
	r.Log(level, args...)

	logger.releaseRecord(r)
}

// Log an message
func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	r := logger.newRecord()
	r.Logf(level, format, args...)

	logger.releaseRecord(r)
}

// Warning logs a message at level Warn
func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

// Warn logs a message at level Warn
func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

// Warnf logs a message at level Warn
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

// Info logs a message at level Info
func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

// Info logs a message at level Info
func (logger *Logger) Infof(format string, args ...interface{})  {
	logger.Logf(InfoLevel, format, args...)
}

// Trace logs a message at level Trace
func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

// Trace logs a message at level Trace
func (logger *Logger) Tracef(format string, args ...interface{})  {
	logger.Logf(TraceLevel, format, args...)
}

// Error logs a message at level Error
func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

// Error logs a message at level Error
func (logger *Logger) Errorf(format string, args ...interface{})  {
	logger.Logf(ErrorLevel, format, args...)
}

// Notice logs a message at level Notice
func (logger *Logger) Notice(args ...interface{}) {
	logger.Log(NoticeLevel, args...)
}

// Notice logs a message at level Notice
func (logger *Logger) Noticef(format string, args ...interface{})  {
	logger.Logf(NoticeLevel, format, args...)
}

// Debug logs a message at level Debug
func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

// Debug logs a message at level Debug
func (logger *Logger) Debugf(format string, args ...interface{})  {
	logger.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
}

// Fatal logs a message at level Fatal
func (logger *Logger) Fatalf(format string, args ...interface{})  {
	logger.Logf(FatalLevel, format, args...)
}

// Panic logs a message at level Panic
func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

// Panic logs a message at level Panic
func (logger *Logger) Panicf(format string, args ...interface{})  {
	logger.Logf(PanicLevel, format, args...)
}

//

// TODO use Record or *Record ...
func (logger *Logger) write(level Level, r *Record) {
	var matchedHandlers []Handler
	for _, handler := range logger.handlers {
		if handler.IsHandling(level) {
			matchedHandlers = append(matchedHandlers, handler)
		}
	}

	// log level is don't match
	if len(matchedHandlers) == 0 {
		return
	}

	if logger.ReportCaller {
		logger.mu.Lock()
		r.Caller = getCaller(logger.MaxCallerDepth)
		logger.mu.Unlock()
	}

	// processing log record
	for _, processor := range logger.processors {
		processor.Process(r)
	}

	// handling log record
	for _, handler := range matchedHandlers {
		err := handler.Handle(r)
		if err != nil {
			_,_ = fmt.Fprintf(os.Stderr, "Failed to dispatch handler: %v\n", err)
			return
		}
	}

	// If is Panic level
	if level <= PanicLevel {
		panic(r)
	}
}
