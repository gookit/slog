package processor

import (
	"os"

	"github.com/gookit/slog"
)

// AddHostname to record
func AddHostname() slog.Processor {
	hostname,_ := os.Hostname()

	return slog.ProcessorFunc(func(record *slog.Record) {
		record.AddField("hostname", hostname)
	})
}
