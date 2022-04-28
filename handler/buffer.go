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

	return &FlushCloseHandler{
		Output: bufwrite.NewBufIOWriterSize(w, bufSize),
		// log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// LineBufferedFile handler
func LineBufferedFile(logfile string, bufSize int, levels []slog.Level) (slog.Handler, error) {
	cfg := NewConfig(
		WithLogfile(logfile),
		WithBuffSize(bufSize),
		WithLogLevels(levels),
		WithBuffMode(BuffModeLine),
	)

	output, err := cfg.CreateWriter()
	if err != nil {
		return nil, err
	}

	return &SyncCloseHandler{
		Output: output,
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(cfg.Levels),
	}, nil
}

// LineBuffOsFile handler
func LineBuffOsFile(f *os.File, bufSize int, levels []slog.Level) slog.Handler {
	if f == nil {
		panic("slog: the os file cannot be nil")
	}

	return &SyncCloseHandler{
		Output: bufwrite.NewLineWriterSize(f, bufSize),
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}

// LineBuffWriter handler
func LineBuffWriter(w io.Writer, bufSize int, levels []slog.Level) slog.Handler {
	if w == nil {
		panic("slog: the io writer cannot be nil")
	}

	return &IOWriterHandler{
		Output: bufwrite.NewLineWriterSize(w, bufSize),
		// init log levels
		LevelFormattable: slog.NewLvsFormatter(levels),
	}
}
