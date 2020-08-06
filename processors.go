package slog

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"runtime"

	"github.com/gookit/goutil/strutil"
)

// there are some built-in processors

// AddHostname to record
func AddHostname() Processor {
	hostname, _ := os.Hostname()

	return ProcessorFunc(func(record *Record) {
		record.AddField("hostname", hostname)
	})
}

// AddUniqueID to record
func AddUniqueID(fieldName string) Processor {
	h := md5.New()

	return ProcessorFunc(func(record *Record) {
		rb, err := strutil.RandomBytes(32)
		if err != nil {
			record.WithError(err)
			return
		}

		h.Write(rb)
		randomId := hex.EncodeToString(h.Sum(nil))
		h.Reset()

		record.AddField(fieldName, randomId)
	})
}

// MemoryUsage Get memory usage.
var MemoryUsage ProcessorFunc = func(record *Record) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	record.Extra["memoryUsage"] = stat.Alloc
}
