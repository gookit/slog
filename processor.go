package slog

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"runtime"

	"github.com/gookit/goutil/strutil"
)

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
	// ProcessRecord handle an record
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

// ProcessRecord process record
func (p *Processable) ProcessRecord(r *Record) {
	// processing log record
	for _, processor := range p.processors {
		processor.Process(r)
	}
}

//
// there are some built-in processors
//

// AddHostname to record
func AddHostname() Processor {
	hostname, _ := os.Hostname()

	return ProcessorFunc(func(record *Record) {
		record.AddField("hostname", hostname)
	})
}

// AddUniqueID to record
func AddUniqueID(fieldName string) Processor {
	hs := md5.New()

	return ProcessorFunc(func(record *Record) {
		rb, err := strutil.RandomBytes(32)
		if err != nil {
			record.WithError(err)
			return
		}

		hs.Write(rb)
		randomID := hex.EncodeToString(hs.Sum(nil))
		hs.Reset()

		record.AddField(fieldName, randomID)
	})
}

// MemoryUsage get memory usage.
var MemoryUsage ProcessorFunc = func(record *Record) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	record.SetExtraValue("memoryUsage", stat.Alloc)
}
