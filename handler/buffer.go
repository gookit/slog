package handler

import (
	"io"
	"os"

	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
)

// NewBuffered create new BufferedHandler
func NewBuffered(w io.WriteCloser, bufSize int, levels ...slog.Level) *FlushCloseHandler {
	return NewBufferedHandler(w, bufSize, levels...)
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(w io.WriteCloser, bufSize int, levels ...slog.Level) *FlushCloseHandler {
	if len(levels) == 0 {
		levels = slog.AllLevels
	}

	out := bufwrite.NewBufIOWriterSize(w, bufSize)
	return FlushCloserWithLevels(out, levels)
}

// LineBufferedFile handler
func LineBufferedFile(logfile string, bufSize int, levels []slog.Level) (slog.Handler, error) {
	cfg := NewConfig(
		WithLogfile(logfile),
		WithBuffSize(bufSize),
		WithLogLevels(levels),
		WithBuffMode(BuffModeLine),
	)

	out, err := cfg.CreateWriter()
	if err != nil {
		return nil, err
	}
	return SyncCloserWithLevels(out, levels), nil
}

// LineBuffOsFile handler
func LineBuffOsFile(f *os.File, bufSize int, levels []slog.Level) slog.Handler {
	if f == nil {
		panic("slog: the os file cannot be nil")
	}

	out := bufwrite.NewLineWriterSize(f, bufSize)
	return SyncCloserWithLevels(out, levels)
}

// LineBuffWriter handler
func LineBuffWriter(w io.Writer, bufSize int, levels []slog.Level) slog.Handler {
	if w == nil {
		panic("slog: the io writer cannot be nil")
	}

	out := bufwrite.NewLineWriterSize(w, bufSize)
	return IOWriterWithLevels(out, levels)
}

//
// --------- wrap a handler with buffer ---------
//

// FormatWriterHandler interface
type FormatWriterHandler interface {
	slog.Handler
	// Formatter record formatter
	Formatter() slog.Formatter
	// Writer the output writer
	Writer() io.Writer
}
