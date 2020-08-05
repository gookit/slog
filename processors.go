package slog

import (
	"os"
	"runtime"
)

// AddHostname to record
func AddHostname() Processor {
	hostname, _ := os.Hostname()

	return ProcessorFunc(func(record *Record) {
		record.AddField("hostname", hostname)
	})
}

// MemoryUsage Get memory usage.
var MemoryUsage ProcessorFunc = func(record *Record) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	record.Extra["memoryUsage"] = stat.Alloc
}
