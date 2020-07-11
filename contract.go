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
	Handle(*Record) error
	// HandleBatch Handles a set of records at once.
	HandleBatch([]*Record) error
}

// GroupedHandler definition
type GroupedHandler struct {
	handlers []Handler
}

//
// processor
//

// Processor interface
type Processor interface {
	Process(record *Record)
}

type ProcessorFunc func(record *Record)

func (fn ProcessorFunc) Process(record *Record)  {
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

func (p *Processable) AddProcessor(processor Processor) {
	p.processors = append(p.processors, processor)
}

func (p *Processable) ProcessRecord(r *Record) {
	// processing log record
	for _, processor := range p.processors {
		processor.Process(r)
	}
}

//
// formatter
//

// Formatter interface
type Formatter interface {
	Format(record *Record)  ([]byte, error)
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

func (f *Formattable) Formatter() Formatter {
	return f.formatter
}

func (f *Formattable) SetFormatter(formatter Formatter) {
	f.formatter = formatter
}

