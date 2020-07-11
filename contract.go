package slog

import (
	"io"
)

// Handler interface
type Handler interface {
	io.Closer
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// all records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) bool
	// HandleBatch Handles a set of records at once.
	HandleBatch([]*Record) bool
}

// Formatter interface
type Formatter interface {
	Format(record *Record)  ([]byte, error)
}

// Processor interface
type Processor interface {
	Process(record *Record)
}

// ProcessableHandler interface
type ProcessableHandler interface {
	// Processor get the log processor
	Processor() Processor
	// SetProcessor set the log processor
	SetProcessor(Processor)
}

// FormattableHandler interface
type FormattableHandler interface {
	// Processor get the log processor
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}

type ProcessorFunc func(record *Record)

func (fn ProcessorFunc) Process(record *Record)  {
	fn(record)
}
