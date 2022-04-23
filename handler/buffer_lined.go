package handler

import (
	"io"
	"os"

	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
)

// LineBufferedFile handler
func LineBufferedFile(fPath string, bufSize int, levels []slog.Level) slog.Handler {
	f, err := QuickOpenFile(fPath)
	if err != nil {
		panic(err)
	}

	return &IOWriterHandler{
		Output: bufwrite.NewLineWriterSize(f, bufSize),
		// init log levels
		LevelsWithFormatter: newLvsFormatter(levels),
	}
}

// LineBuffOsFile handler
func LineBuffOsFile(f *os.File, bufSize int, levels []slog.Level) slog.Handler {
	if f == nil {
		panic("slog: the os file cannot be nil")
	}

	return &IOWriterHandler{
		Output: bufwrite.NewLineWriterSize(f, bufSize),
		// init log levels
		LevelsWithFormatter: newLvsFormatter(levels),
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
		LevelsWithFormatter: newLvsFormatter(levels),
	}
}
