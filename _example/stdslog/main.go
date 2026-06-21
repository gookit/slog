// Demo: use gookit/slog's rotatefile.Writer as the destination for the standard
// library log/slog (Go 1.21+). rotatefile.Writer is a plain io.Writer, so it
// plugs straight into a slog handler — you get file rotation + cleanup + gzip
// while keeping the standard slog API.
//
// Run: go run ./_example/stdslog
package main

import (
	"log/slog"

	"github.com/gookit/rotatefile"
)

func main() {
	// rotate by size/time, keep N backups, gzip — all from rotatefile config
	w, err := rotatefile.NewConfig("testdata/std_slog.log", func(c *rotatefile.Config) {
		c.MaxSize = 50 * 1024 * 1024 // 50MB
		c.RotateTime = rotatefile.EveryDay
		c.BackupNum = 7
		c.Compress = true
	}).Create()
	if err != nil {
		panic(err)
	}
	defer w.Close() // flush + close on exit

	logger := slog.New(slog.NewJSONHandler(w, nil))
	logger.Info("log message via std log/slog", "app", "demo", "n", 1)
	logger.Warn("a warning", slog.Group("req", "method", "GET", "path", "/ping"))
}
