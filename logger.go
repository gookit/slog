package slog

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

// M short name of map[string]interface{}
type M map[string]interface{}

// Logger definition
type Logger struct {
	// name string

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
	ExitFunc     func(int)
}

// New create an new logger
func New() *Logger {
	logger := &Logger{
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
	logger.recordPool.Put(r)
}

// Register handlers and processors

// AddHandler to the logger
func (logger *Logger) AddHandler(h Handler) {
	logger.handlers = append(logger.handlers, h)
}

// AddHandlers to the logger
func (logger *Logger) AddHandlers(hs ...Handler) {
	logger.handlers = append(logger.handlers, hs...)
}

// PushHandler to the logger. alias of PushHandler()
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
func (logger *Logger) AddProcessors(ps []Processor) {
	logger.processors = append(logger.processors, ps...)
}

// SetProcessors for the logger
func (logger *Logger) SetProcessors(ps []Processor) {
	logger.processors = ps
}

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

//

// TODO use Record or *Record ...
func (logger *Logger) write(level Level, r Record) {
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
		processor.Process(&r)
	}

	// handling log record
	for _, handler := range matchedHandlers {
		err := handler.Handle(&r)
		if err != nil {
			_,_ = fmt.Fprintf(os.Stderr, "Failed to dispatch handler: %v\n", err)
			return
		}
	}

	var buffer *bytes.Buffer

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	// If is Panic level
	if level <= PanicLevel {
		panic(&r)
	}
}
