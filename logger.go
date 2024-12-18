package slog

import (
	"context"
	"sync"
	"time"

	"github.com/gookit/goutil"
)

// Logger log dispatcher definition.
//
// The logger implements the `github.com/gookit/gsr.Logger`
type Logger struct {
	name string
	// lock for write logs
	mu sync.Mutex
	// logger latest error
	err error
	// mark logger is closed
	closed bool

	// log handlers for logger
	handlers   []Handler
	processors []Processor

	// reusable empty record
	recordPool sync.Pool
	// handlers on exit.
	exitHandlers []func()
	quitDaemon   chan struct{}

	//
	// logger options
	//

	// ChannelName log channel name, default is DefaultChannelName
	ChannelName string
	// FlushInterval flush interval time. default is defaultFlushInterval=30s
	FlushInterval time.Duration
	// LowerLevelName use lower level name
	LowerLevelName bool
	// ReportCaller on write log record
	ReportCaller bool
	CallerSkip   int
	// CallerFlag used to set caller traceback information in different modes
	CallerFlag CallerFlagMode
	// BackupArgs backup log input args to Record.Args
	BackupArgs bool
	// TimeClock custom time clock, timezone
	TimeClock ClockFn
	// custom exit, panic handler.
	ExitFunc  func(code int)
	PanicFunc func(v any)
}

// New create a new logger
func New(fns ...LoggerFn) *Logger {
	return NewWithName("logger", fns...)
}

// NewWithHandlers create a new logger with handlers
func NewWithHandlers(hs ...Handler) *Logger {
	logger := NewWithName("logger")
	logger.AddHandlers(hs...)
	return logger
}

// NewWithConfig create a new logger with config func
func NewWithConfig(fns ...LoggerFn) *Logger {
	return NewWithName("logger", fns...)
}

// NewWithName create a new logger with name
func NewWithName(name string, fns ...LoggerFn) *Logger {
	logger := &Logger{
		name: name,
		// exit handle
		// ExitFunc:  os.Exit,
		PanicFunc:    DefaultPanicFn,
		exitHandlers: []func(){},
		// options
		ChannelName:  DefaultChannelName,
		ReportCaller: true,
		CallerSkip:   6,
		TimeClock:    DefaultClockFn,
		// flush interval time
		FlushInterval: defaultFlushInterval,
	}

	logger.recordPool.New = func() any {
		return newRecord(logger)
	}
	return logger.Config(fns...)
}

// NewRecord get new logger record
func (l *Logger) newRecord() *Record {
	r := l.recordPool.Get().(*Record)
	r.reuse = false
	r.freed = false
	r.Fields = nil
	return r
}

func (l *Logger) releaseRecord(r *Record) {
	if r.reuse || r.freed {
		return
	}

	// reset data
	r.Time = emptyTime
	r.Data = map[string]any{}
	r.Extra = nil
	// reset flags
	r.inited = false
	r.reuse = false
	r.freed = true

	r.Message = ""
	r.CallerSkip = l.CallerSkip
	l.recordPool.Put(r)
}

//
// ---------------------------------------------------------------------------
// Configure logger
// ---------------------------------------------------------------------------
//

// Config current logger
func (l *Logger) Config(fns ...LoggerFn) *Logger {
	for _, fn := range fns {
		fn(l)
	}
	return l
}

// Configure current logger. alias of Config()
func (l *Logger) Configure(fn LoggerFn) *Logger { return l.Config(fn) }

// RegisterExitHandler register an exit-handler on global exitHandlers
func (l *Logger) RegisterExitHandler(handler func()) {
	l.exitHandlers = append(l.exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func (l *Logger) PrependExitHandler(handler func()) {
	l.exitHandlers = append([]func(){handler}, l.exitHandlers...)
}

// ResetExitHandlers reset logger exitHandlers
func (l *Logger) ResetExitHandlers() {
	l.exitHandlers = make([]func(), 0)
}

// ExitHandlers get all exitHandlers of the logger
func (l *Logger) ExitHandlers() []func() {
	return l.exitHandlers
}

// SetName for logger
func (l *Logger) SetName(name string) { l.name = name }

// Name of the logger
func (l *Logger) Name() string { return l.name }

//
// ---------------------------------------------------------------------------
// Management logger
// ---------------------------------------------------------------------------
//

const defaultFlushInterval = 30 * time.Second

// FlushDaemon run flush handle on daemon
//
// Usage please refer to the FlushDaemon() on package.
func (l *Logger) FlushDaemon(onStops ...func()) {
	l.quitDaemon = make(chan struct{})
	if l.FlushInterval <= 0 {
		l.FlushInterval = defaultFlushInterval
	}

	// create a ticker
	tk := time.NewTicker(l.FlushInterval)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			if err := l.lockAndFlushAll(); err != nil {
				printlnStderr("slog.FlushDaemon: daemon flush logs error: ", err)
			}
		case <-l.quitDaemon:
			for _, fn := range onStops {
				fn()
			}
			return
		}
	}
}

// StopDaemon stop flush daemon
func (l *Logger) StopDaemon() {
	if l.quitDaemon == nil {
		panic("cannot quit daemon, please call FlushDaemon() first")
	}
	close(l.quitDaemon)
}

// FlushTimeout flush logs on limit time.
//
// refer from glog package
func (l *Logger) FlushTimeout(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		if err := l.lockAndFlushAll(); err != nil {
			printlnStderr("slog.FlushTimeout: flush logs error: ", err)
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		printlnStderr("slog.FlushTimeout: flush took longer than timeout:", timeout)
	}
}

// Sync flushes buffered logs (if any). alias of the Flush()
func (l *Logger) Sync() error { return Flush() }

// Flush flushes all the logs and attempts to "sync" their data to disk.
// l.mu is held.
func (l *Logger) Flush() error { return l.lockAndFlushAll() }

// MustFlush flush logs. will panic on error
func (l *Logger) MustFlush() {
	goutil.PanicErr(l.lockAndFlushAll())
}

// FlushAll flushes all the logs and attempts to "sync" their data to disk.
//
// alias of the Flush()
func (l *Logger) FlushAll() error { return l.lockAndFlushAll() }

// lockAndFlushAll is like flushAll but locks l.mu first.
func (l *Logger) lockAndFlushAll() error {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()

	return l.err
}

// flush all without lock
func (l *Logger) flushAll() {
	// flush from fatal down, in case there's trouble flushing.
	_ = l.VisitAll(func(handler Handler) error {
		if err := handler.Flush(); err != nil {
			l.err = err
			printlnStderr("slog: call handler.Flush() error:", err)
		}
		return nil
	})
}

// MustClose close logger. will panic on error
func (l *Logger) MustClose() {
	goutil.PanicErr(l.Close())
}

// Close the logger, will flush all logs and close all handlers
//
// IMPORTANT:
//
//	if enable async/buffer mode, please call the Close() before exit.
func (l *Logger) Close() error {
	if l.closed {
		return nil
	}

	_ = l.VisitAll(func(handler Handler) error {
		if err := handler.Close(); err != nil {
			l.err = err
			printlnStderr("slog: call handler.Close() error:", err)
		}
		return nil
	})

	l.closed = true
	return l.err
}

// VisitAll logger handlers
func (l *Logger) VisitAll(fn func(handler Handler) error) error {
	for _, handler := range l.handlers {
		// TIP: you can return nil for ignore error
		if err := fn(handler); err != nil {
			return err
		}
	}
	return nil
}

// Reset the logger. will reset: handlers, processors, closed=false
func (l *Logger) Reset() {
	l.closed = false
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

	if l.ExitFunc != nil {
		l.ExitFunc(code)
	}
}

func (l *Logger) runExitHandlers() {
	defer func() {
		if err := recover(); err != nil {
			printlnStderr("slog: run exit handler recovered, error:", err)
		}
	}()

	for _, handler := range l.exitHandlers {
		handler()
	}
}

// DoNothingOnPanicFatal do nothing on panic or fatal level. useful on testing.
func (l *Logger) DoNothingOnPanicFatal() {
	l.PanicFunc = DoNothingOnPanic
	l.ExitFunc = DoNothingOnExit
}

// HandlersNum returns the number of handlers
func (l *Logger) HandlersNum() int {
	return len(l.handlers)
}

// LastErr fetch, will clear it after read.
func (l *Logger) LastErr() error {
	err := l.err
	l.err = nil
	return err
}

//
// ---------------------------------------------------------------------------
// Register handlers and processors
// ---------------------------------------------------------------------------
//

// AddHandler to the logger
func (l *Logger) AddHandler(h Handler) { l.PushHandlers(h) }

// AddHandlers to the logger
func (l *Logger) AddHandlers(hs ...Handler) { l.PushHandlers(hs...) }

// PushHandler to the l. alias of AddHandler()
func (l *Logger) PushHandler(h Handler) { l.PushHandlers(h) }

// PushHandlers to the logger
func (l *Logger) PushHandlers(hs ...Handler) {
	if len(hs) > 0 {
		l.handlers = append(l.handlers, hs...)
	}
}

// SetHandlers for the logger
func (l *Logger) SetHandlers(hs []Handler) { l.handlers = hs }

// AddProcessor to the logger
func (l *Logger) AddProcessor(p Processor) { l.processors = append(l.processors, p) }

// PushProcessor to the logger, alias of AddProcessor()
func (l *Logger) PushProcessor(p Processor) { l.processors = append(l.processors, p) }

// AddProcessors to the logger. alias of AddProcessor()
func (l *Logger) AddProcessors(ps ...Processor) { l.processors = append(l.processors, ps...) }

// SetProcessors for the logger
func (l *Logger) SetProcessors(ps []Processor) { l.processors = ps }

//
// ---------------------------------------------------------------------------
// New record with log data, fields
// ---------------------------------------------------------------------------
//

// Record return a new record with logger, will release after write log.
func (l *Logger) Record() *Record {
	return l.newRecord()
}

// Reused return a new record with logger, but it can be reused.
// if you want to release the record, please call the Record.Release() after write log.
//
// Usage:
//
//	r := logger.Reused()
//	defer r.Release()
//
//	// can write log multiple times
//	r.Info("some message1")
//	r.Warn("some message1")
func (l *Logger) Reused() *Record {
	return l.newRecord().Reused()
}

// WithField new record with field
//
// TIP: add field need config Formatter template fields.
func (l *Logger) WithField(name string, value any) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)
	return r.WithField(name, value)
}

// WithFields new record with fields
//
// TIP: add field need config Formatter template fields.
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

// WithValue new record with data value
func (l *Logger) WithValue(key string, value any) *Record {
	r := l.newRecord()
	return r.AddValue(key, value)
}

// WithExtra new record with extra data
func (l *Logger) WithExtra(ext M) *Record {
	r := l.newRecord()
	return r.SetExtra(ext)
}

// WithTime new record with time.Time
func (l *Logger) WithTime(t time.Time) *Record {
	r := l.newRecord()
	defer l.releaseRecord(r)
	return r.WithTime(t)
}

// WithCtx new record with context.Context
func (l *Logger) WithCtx(ctx context.Context) *Record { return l.WithContext(ctx) }

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

func (l *Logger) log(level Level, args []any) {
	r := l.newRecord()
	r.CallerSkip++
	r.log(level, args)
}

// Logf a format message with level
func (l *Logger) logf(level Level, format string, args []any) {
	r := l.newRecord()
	r.CallerSkip++
	r.logf(level, format, args)
}

// Log a message with level
func (l *Logger) Log(level Level, args ...any) { l.log(level, args) }

// Logf a format message with level
func (l *Logger) Logf(level Level, format string, args ...any) {
	l.logf(level, format, args)
}

// Print logs a message at level PrintLevel
func (l *Logger) Print(args ...any) { l.log(PrintLevel, args) }

// Println logs a message at level PrintLevel
func (l *Logger) Println(args ...any) { l.log(PrintLevel, args) }

// Printf logs a message at level PrintLevel
func (l *Logger) Printf(format string, args ...any) { l.logf(PrintLevel, format, args) }

// Warn logs a message at level Warn
func (l *Logger) Warn(args ...any) { l.log(WarnLevel, args) }

// Warnf logs a message at level Warn
func (l *Logger) Warnf(format string, args ...any) { l.logf(WarnLevel, format, args) }

// Warning logs a message at level Warn, alias of Logger.Warn()
func (l *Logger) Warning(args ...any) { l.log(WarnLevel, args) }

// Info logs a message at level Info
func (l *Logger) Info(args ...any) { l.log(InfoLevel, args) }

// Infof logs a message at level Info
func (l *Logger) Infof(format string, args ...any) { l.logf(InfoLevel, format, args) }

// Trace logs a message at level trace
func (l *Logger) Trace(args ...any) { l.log(TraceLevel, args) }

// Tracef logs a message at level trace
func (l *Logger) Tracef(format string, args ...any) { l.logf(TraceLevel, format, args) }

// Error logs a message at level error
func (l *Logger) Error(args ...any) { l.log(ErrorLevel, args) }

// Errorf logs a message at level error
func (l *Logger) Errorf(format string, args ...any) { l.logf(ErrorLevel, format, args) }

// ErrorT logs a error type at level error
func (l *Logger) ErrorT(err error) {
	if err != nil {
		l.log(ErrorLevel, []any{err})
	}
}

// Stack logs a error message and with call stack. TODO
// func EStack(args ...any) { std.log(ErrorLevel, args) }

// Notice logs a message at level notice
func (l *Logger) Notice(args ...any) { l.log(NoticeLevel, args) }

// Noticef logs a message at level notice
func (l *Logger) Noticef(format string, args ...any) { l.logf(NoticeLevel, format, args) }

// Debug logs a message at level debug
func (l *Logger) Debug(args ...any) { l.log(DebugLevel, args) }

// Debugf logs a message at level debug
func (l *Logger) Debugf(format string, args ...any) { l.logf(DebugLevel, format, args) }

// Fatal logs a message at level fatal
func (l *Logger) Fatal(args ...any) { l.log(FatalLevel, args) }

// Fatalf logs a message at level fatal
func (l *Logger) Fatalf(format string, args ...any) { l.logf(FatalLevel, format, args) }

// Fatalln logs a message at level fatal
func (l *Logger) Fatalln(args ...any) { l.log(FatalLevel, args) }

// Panic logs a message at level panic
func (l *Logger) Panic(args ...any) { l.log(PanicLevel, args) }

// Panicf logs a message at level panic
func (l *Logger) Panicf(format string, args ...any) { l.logf(PanicLevel, format, args) }

// Panicln logs a message at level panic
func (l *Logger) Panicln(args ...any) { l.log(PanicLevel, args) }
