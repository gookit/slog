package handler

import (
	"io"

	"github.com/gookit/slog"
	"github.com/gookit/slog/rotatefile"
)

//
// ---------------------------------------------------------------------------
// handler builder
// ---------------------------------------------------------------------------
//

// Builder struct for create handler
type Builder struct {
	*Config
	Output io.Writer
}

// NewBuilder create
func NewBuilder() *Builder {
	return &Builder{
		Config: NewEmptyConfig(),
	}
}

// WithOutput to the builder
func (b *Builder) WithOutput(w io.Writer) *Builder {
	b.Output = w
	return b
}

// With some config fn
//
// Deprecated: please use WithConfigFn()
func (b *Builder) With(fns ...ConfigFn) *Builder {
	return b.WithConfigFn(fns...)
}

// WithConfigFn some config fn
func (b *Builder) WithConfigFn(fns ...ConfigFn) *Builder {
	b.Config.With(fns...)
	return b
}

// WithLogfile setting
func (b *Builder) WithLogfile(logfile string) *Builder {
	b.Logfile = logfile
	return b
}

// WithLevelMode setting
func (b *Builder) WithLevelMode(mode uint8) *Builder {
	b.LevelMode = mode
	return b
}

// WithLogLevel setting
func (b *Builder) WithLogLevel(level slog.Level) *Builder {
	b.Level = level
	return b
}

// WithLogLevels setting
func (b *Builder) WithLogLevels(levels []slog.Level) *Builder {
	b.Levels = levels
	return b
}

// WithBuffMode setting
func (b *Builder) WithBuffMode(bufMode string) *Builder {
	b.BuffMode = bufMode
	return b
}

// WithBuffSize setting
func (b *Builder) WithBuffSize(bufSize int) *Builder {
	b.BuffSize = bufSize
	return b
}

// WithMaxSize setting
func (b *Builder) WithMaxSize(maxSize uint64) *Builder {
	b.MaxSize = maxSize
	return b
}

// WithRotateTime setting
func (b *Builder) WithRotateTime(rt rotatefile.RotateTime) *Builder {
	b.RotateTime = rt
	return b
}

// WithCompress setting
func (b *Builder) WithCompress(compress bool) *Builder {
	b.Compress = compress
	return b
}

// WithUseJSON setting
func (b *Builder) WithUseJSON(useJSON bool) *Builder {
	b.UseJSON = useJSON
	return b
}

// Build slog handler.
func (b *Builder) Build() slog.FormattableHandler {
	if b.Output != nil {
		return b.buildFromWriter(b.Output)
	}

	if b.Logfile != "" {
		w, err := b.CreateWriter()
		if err != nil {
			panic(err)
		}
		return b.buildFromWriter(w)
	}

	panic("slog: missing information for build slog handler")
}

// Build slog handler.
func (b *Builder) buildFromWriter(w io.Writer) (h slog.FormattableHandler) {
	defer b.reset()
	bufSize := b.BuffSize
	lf := b.newLevelFormattable()

	if scw, ok := w.(SyncCloseWriter); ok {
		if bufSize > 0 {
			scw = b.wrapBuffer(scw)
		}

		h = &SyncCloseHandler{
			Output: scw,
			// with log level and formatter
			LevelFormattable: lf,
		}
	} else if fcw, ok := w.(FlushCloseWriter); ok {
		if bufSize > 0 {
			fcw = b.wrapBuffer(fcw)
		}

		h = &FlushCloseHandler{
			Output: fcw,
			// with log level and formatter
			LevelFormattable: lf,
		}
	} else if wc, ok := w.(io.WriteCloser); ok {
		if bufSize > 0 {
			wc = b.wrapBuffer(wc)
		}

		h = &WriteCloserHandler{
			Output: wc,
			// with log level and formatter
			LevelFormattable: lf,
		}
	} else {
		if bufSize > 0 {
			w = b.wrapBuffer(w)
		}

		h = &IOWriterHandler{
			Output: w,
			// with log level and formatter
			LevelFormattable: lf,
		}
	}

	// use json format.
	if b.UseJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	}
	return
}

// rest builder.
func (b *Builder) reset() {
	b.Output = nil
	b.Config = NewEmptyConfig()
}
