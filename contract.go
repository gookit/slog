package slog

import (
	"io"
)

// SyncWriter is the interface satisfied by logging destinations.
type SyncWriter interface {
	Sync() error
	io.Writer
}

// FlushSyncWriter is the interface satisfied by logging destinations.
type FlushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

// WriterHandler is the interface satisfied by logging destinations.
type WriterHandler interface {
	Handler
	Writer() io.Writer
}

//
// Handler interface
//

// Handler interface definition
type Handler interface {
	// Close handler
	io.Closer
	// Flush logs to disk
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	// all records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
	// HandleBatch Handles a set of records at once. TODO need ?
	// HandleBatch([]*Record) error
}

//
// Processor interface
//

// Processor interface definition
type Processor interface {
	// Process record
	Process(record *Record)
}

// ProcessorFunc wrapper definition
type ProcessorFunc func(record *Record)

// Process record
func (fn ProcessorFunc) Process(record *Record) {
	fn(record)
}

// ProcessableHandler interface
type ProcessableHandler interface {
	// AddProcessor add an processor
	AddProcessor(Processor)
	// SetProcessor set the log processor
	ProcessRecord(record *Record)
}

// Processable definition
type Processable struct {
	processors []Processor
}

// AddProcessor to the handler
func (p *Processable) AddProcessor(processor Processor) {
	p.processors = append(p.processors, processor)
}

// ProcessRecord process records
func (p *Processable) ProcessRecord(r *Record) {
	// processing log record
	for _, processor := range p.processors {
		processor.Process(r)
	}
}

//
// Formatter interface
//

// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}

// FormatterFunc wrapper definition
type FormatterFunc func(r *Record) ([]byte, error)

// Format an record
func (fn FormatterFunc) Format(r *Record) ([]byte, error) {
	return fn(r)
}

// FormattableHandler interface
type FormattableHandler interface {
	// Formatter get the log formatter
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}

// Formattable definition
type Formattable struct {
	formatter Formatter
}

// Formatter get formatter. if not set, will return TextFormatter
func (f *Formattable) Formatter() Formatter {
	if f.formatter == nil {
		f.formatter = NewTextFormatter()
	}
	return f.formatter
}

// SetFormatter to handler
func (f *Formattable) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}
