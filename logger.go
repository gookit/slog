package slog

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger definition.
// The logger implements the `github.com/gookit/gsr.Logger`
type Logger struct {
	name string

	// timezone
	tz time.Time
	mu sync.Mutex

	handlers   []Handler
	processors []Processor

	// options
	// ReportCaller on log message
	ReportCaller   bool
	LowerLevelName bool
	MaxCallerDepth int

	// Reusable empty record
	recordPool sync.Pool

	// handlers on exit
	exitHandlers []func()
	ExitFunc     func(code int)
}

// New create an new logger
func New() *Logger {
	return NewWithName("-")
}

// NewWithConfig create an new logger with config func
func NewWithConfig(fn func(l *Logger)) *Logger {
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
		ReportCaller:   true,
		MaxCallerDepth: defaultMaxCallerDepth,
	}

	logger.recordPool.New = func() interface{} {
		return newRecord(logger)
	}

	return logger
}

// NewRecord get new logger record
func (l *Logger) newRecord() *Record {
	return l.recordPool.Get().(*Record)
}

func (l *Logger) releaseRecord(r *Record) {
	// reset data
	r.Time = time.Time{}
	r.Data = nil
	r.Extra = nil
	r.Fields = nil
	l.recordPool.Put(r)
}

//
// ---------------------------------------------------------------------------
// Management logger
// ---------------------------------------------------------------------------
//

const flushInterval = 30 * time.Second

// Configure current logger
func (l *Logger) Configure(fn func(l *Logger)) *Logger {
	fn(l)
	return l
}

// Sync flushes buffered logs (if any).
func (l *Logger) Sync() error {
	// TODO
	return nil
}

// FlushDaemon run flush handle on daemon
//
// Usage:
// 	go slog.FlushDaemon()
func (l *Logger) FlushDaemon() {
	for range time.NewTicker(flushInterval).C {
		l.lockAndFlushAll()
	}
}

// FlushTimeout flush logs on limit time.
// refer from glog package
func (l *Logger) FlushTimeout(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		l.FlushAll() // calls l.lockAndFlushAll()
		// _, _ = fmt.Fprintln(os.Stderr, "slog: Flush logs error.", err)

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		_, _ = fmt.Fprintln(os.Stderr, "slog: Flush took longer than", timeout)
	}
}

// lockAndFlushAll is like flushAll but locks l.mu first.
func (l *Logger) lockAndFlushAll() {
	l.mu.Lock()
	l.FlushAll()
	l.mu.Unlock()
}

// Flush flushes all the logs to disk. alias of the FlushAll()
func (l *Logger) Flush() {
	l.FlushAll()
}

// FlushAll flushes all the logs and attempts to "sync" their data to disk.
// l.mu is held.
func (l *Logger) FlushAll() {
	// Flush from fatal down, in case there's trouble flushing.
	l.VisitAll(func(handler Handler) error {
		_ = handler.Flush() // ignore error
		return nil
	})
}

// Close the logger
func (l *Logger) Close() {
	l.VisitAll(func(handler Handler) error {
		// Flush logs and then close
		_ = handler.Close() // ignore error
		return nil
	})
}

// VisitAll logger handlers
func (l *Logger) VisitAll(fn func(handler Handler) error) {
	for _, handler := range l.handlers {
		// you can return nil for ignore error
		if err := fn(handler); err != nil {
			return
		}
	}
}

// Reset the logger
func (l *Logger) Reset() {
	l.ResetHandlers()
	l.ResetProcessors()
}

// ResetProcessors for the logger
func (l *Logger) ResetProcessors() {
	l.processors = make([]Processor, 0)
}

// ResetHandlers for the logger
func (l *Logger) ResetHandlers() {
	l.handlers = make([]Handler, 0)
}

// Exit logger handle
func (l *Logger) Exit(code int) {
	l.runExitHandlers()

	// global exit handlers
	runExitHandlers()

	if l.ExitFunc == nil {
		l.ExitFunc = os.Exit
	}
	l.ExitFunc(code)
}

// SetName for logger
func (l *Logger) SetName(name string) {
	l.name = name
}

// Name of the logger
func (l *Logger) Name() string {
	return l.name
}

//
// ---------------------------------------------------------------------------
// Register handlers and processors
// ---------------------------------------------------------------------------
//

// AddHandler to the logger
func (l *Logger) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

// AddHandlers to the logger
func (l *Logger) AddHandlers(hs ...Handler) {
	l.handlers = append(l.handlers, hs...)
}

// PushHandlers to the logger
func (l *Logger) PushHandlers(hs ...Handler) {
	l.handlers = append(l.handlers, hs...)
}

// PushHandler to the l. alias of AddHandler()
func (l *Logger) PushHandler(h Handler) {
	l.AddHandler(h)
}

// SetHandlers for the logger
func (l *Logger) SetHandlers(hs []Handler) {
	l.handlers = hs
}

// AddProcessor to the logger
func (l *Logger) AddProcessor(p Processor) {
	l.processors = append(l.processors, p)
}

// PushProcessor to the logger
// alias of AddProcessor()
func (l *Logger) PushProcessor(p Processor) {
	l.processors = append(l.processors, p)
}

// AddProcessors to the logger
func (l *Logger) AddProcessors(ps ...Processor) {
	l.processors = append(l.processors, ps...)
}

// SetProcessors for the logger
func (l *Logger) SetProcessors(ps []Processor) {
	l.processors = ps
}

//
// ---------------------------------------------------------------------------
// New record with log data, fields
// ---------------------------------------------------------------------------
//

// SetFields new record with fields
func (l *Logger) WithFields(fields M) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)

	return r.WithFields(fields)
}

// WithData new record with data
func (l *Logger) WithData(data M) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)

	return r.WithData(data)
}

// WithTime new record with time.Time
func (l *Logger) WithTime(t time.Time) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)

	return r.WithTime(t)
}

// WithContext new record with context.Context
func (l *Logger) WithContext(ctx context.Context) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)

	return r.WithContext(ctx)
}

//
// ---------------------------------------------------------------------------
// Add log message with level
// ---------------------------------------------------------------------------
//

// Log an message
func (l *Logger) Log(level Level, args ...interface{}) {
	r := l.newRecord()
	r.Log(level, args...)

	l.releaseRecord(r)
}

// Log an message
func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	r := l.newRecord()
	r.Logf(level, format, args...)

	l.releaseRecord(r)
}

// Print logs a message at level PrintLevel
func (l *Logger) Print(args ...interface{}) {
	l.Log(PrintLevel, args...)
}

// Println logs a message at level PrintLevel
func (l *Logger) Println(args ...interface{}) {
	l.Log(PrintLevel, args...)
}

// Printf logs a message at level PrintLevel
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logf(PrintLevel, format, args...)
}

// Warning logs a message at level Warn
func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

// Warn logs a message at level Warn
func (l *Logger) Warn(args ...interface{}) {
	l.Log(WarnLevel, args...)
}

// Warnf logs a message at level Warn
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(WarnLevel, format, args...)
}

// Info logs a message at level Info
func (l *Logger) Info(args ...interface{}) {
	l.Log(InfoLevel, args...)
}

// Info logs a message at level Info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(InfoLevel, format, args...)
}

// Trace logs a message at level Trace
func (l *Logger) Trace(args ...interface{}) {
	l.Log(TraceLevel, args...)
}

// Trace logs a message at level Trace
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(TraceLevel, format, args...)
}

// Error logs a message at level error
func (l *Logger) Error(args ...interface{}) {
	l.Log(ErrorLevel, args...)
}

// Error logs a message at level Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLevel, format, args...)
}

// ErrorT logs a error type at level Error
func (l *Logger) ErrorT(err error) {
	if err != nil {
		l.Log(ErrorLevel, err)
	}
}

// Notice logs a message at level Notice
func (l *Logger) Notice(args ...interface{}) {
	l.Log(NoticeLevel, args...)
}

// Notice logs a message at level Notice
func (l *Logger) Noticef(format string, args ...interface{}) {
	l.Logf(NoticeLevel, format, args...)
}

// Debug logs a message at level Debug
func (l *Logger) Debug(args ...interface{}) {
	l.Log(DebugLevel, args...)
}

// Debug logs a message at level Debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(FatalLevel, args...)
}

// Fatalf logs a message at level Fatal
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(FatalLevel, format, args...)
}

// Fatalln logs a message at level Fatal
func (l *Logger) Fatalln(args ...interface{}) {
	l.Log(FatalLevel, args...)
}

// Panic logs a message at level Panic
func (l *Logger) Panic(args ...interface{}) {
	l.Log(PanicLevel, args...)
}

// Panicf logs a message at level Panic
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logf(PanicLevel, format, args...)
}

// Panicln logs a message at level Panic
func (l *Logger) Panicln(args ...interface{}) {
	l.Log(PanicLevel, args...)
}

//
// ---------------------------------------------------------------------------
// Do handling log message
// ---------------------------------------------------------------------------
//

func (l *Logger) write(level Level, r *Record) {
	var matchedHandlers []Handler
	for _, handler := range l.handlers {
		if handler.IsHandling(level) {
			matchedHandlers = append(matchedHandlers, handler)
		}
	}

	// log level is don't match
	if len(matchedHandlers) == 0 {
		return
	}

	// use lower level name
	if l.LowerLevelName {
		r.levelName = level.LowerName()
	} else {
		r.levelName = level.Name()
	}

	if l.ReportCaller {
		l.mu.Lock()
		r.Caller = getCaller(l.MaxCallerDepth)
		l.mu.Unlock()
	}

	// processing log record
	for _, processor := range l.processors {
		processor.Process(r)
	}

	// handling log record
	for _, handler := range matchedHandlers {
		if err := handler.Handle(r); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to dispatch handler: %v\n", err)
			return
		}
	}

	// If is Panic level
	if level <= PanicLevel {
		panic(r)
		// If is FatalLevel
	} else if level <= FatalLevel {
		l.Exit(1)
	}
}
