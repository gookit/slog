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

// bufferWrapper struct
type bufferWrapper struct {
	buffer  FlushWriter
	handler FormatWriterHandler
}

// BufferWrapper new instance.
//
// Deprecated: use `NewBufferedHandler` instead, will remove this func at v1.0
func BufferWrapper(h FormatWriterHandler, buffSize int) slog.Handler {
	return &bufferWrapper{
		handler: h,
		buffer:  bufwrite.NewBufIOWriterSize(h.Writer(), buffSize),
	}
}

// IsHandling Check if the current level can be handling
func (w *bufferWrapper) IsHandling(level slog.Level) bool {
	return w.handler.IsHandling(level)
}

// Flush all buffers to the `h.fcWriter.Writer()`
func (w *bufferWrapper) Flush() error {
	if err := w.buffer.Flush(); err != nil {
		return err
	}
	return w.handler.Flush()
}

// Close log records
func (w *bufferWrapper) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.handler.Close()
}

// Handle log record
func (w *bufferWrapper) Handle(record *slog.Record) error {
	bts, err := w.handler.Formatter().Format(record)
	if err != nil {
		return err
	}

	_, err = w.buffer.Write(bts)
	return err
}
