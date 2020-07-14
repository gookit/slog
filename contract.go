package slog

import (
	"io"
)

// M short name of map[string]interface{}
type M map[string]interface{}

// func (m M) String() string  {
// 	return fmt.Sprint(m)
// }

// flushSyncWriter is the interface satisfied by logging destinations.
type FlushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

//
// handler
//

// Handler interface
type Handler interface {
	io.Closer
	Flush() error
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
	// Process record
	Process(record *Record)
}

// ProcessorFunc wrapper definition
type ProcessorFunc func(record *Record)

// Process record
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
// formatter
//

// Formatter interface
type Formatter interface {
	Format(record *Record) ([]byte, error)
}

// FormattableHandler interface
type FormattableHandler interface {
	// Formatter get the log formatter
	Formatter() Formatter
	// SetFormatter set the log formatter
	SetFormatter(Formatter)
}
