package processor

import (
	"os"

	"github.com/gookit/slog"
)

var hostname,_ = os.Hostname()

var AddHostname = slog.ProcessorFunc(func(record *slog.Record) {
	record.AddField("hostname", hostname)
})

