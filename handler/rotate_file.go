package handler

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/gookit/slog"
)

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

// RotateFileHandler definition
type RotateFileHandler struct {
	FileHandler
	logger  *slog.Logger
	written uint64
}

func (h *RotateFileHandler) Write(p []byte) (n int, err error) {
	if h.written+uint64(len(p)) >= h.MaxSize {
		if err := h.rotateFile(time.Now()); err != nil {
			return 0, err
		}
	}

	n, err = h.file.Write(p)
	h.written += uint64(n)
	// if err != nil {
	// 	h.logger.Exit(err)
	// }
	return
}

// -------- refer from glog package
// rotateFile closes the syncBuffer's file and starts a new one.
func (h *RotateFileHandler) rotateFile(now time.Time) error {
	if h.file != nil {
		h.Flush()
		h.file.Close()
	}

	var err error
	h.file, _, err = create("INFO", now)
	h.written = 0
	if err != nil {
		return err
	}

	// init writer
	// h.Writer = bufio.NewWriterSize(h.file, bufferSize)

	// Write header.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
	n, err := h.file.Write(buf.Bytes())

	h.written += uint64(n)
	return err
}
