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
	hs := md5.New()

	return ProcessorFunc(func(record *Record) {
		rb, err := strutil.RandomBytes(32)
		if err != nil {
			record.WithError(err)
			return
		}

		hs.Write(rb)
		randomId := hex.EncodeToString(hs.Sum(nil))
		hs.Reset()

		record.AddField(fieldName, randomId)
	})
}

// MemoryUsage get memory usage.
var MemoryUsage ProcessorFunc = func(record *Record) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	record.SetExtraValue("memoryUsage", stat.Alloc)
}
