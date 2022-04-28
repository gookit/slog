package handler

import (
	"io"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
	"github.com/gookit/slog/rotatefile"
)

const (
	BuffModeLine = "line"
	BuffModeBite = "bite"
)

// ConfigFn for config some settings
type ConfigFn func(c *Config)

// Config struct
type Config struct {
	// Logfile for write logs
	Logfile string `json:"logfile" yaml:"logfile"`
	// UseJSON for format logs
	UseJSON bool `json:"use_json" yaml:"use_json"`
	// BuffMode type name. allow: line, bite
	BuffMode string `json:"buff_mode" yaml:"buff_mode"`
	// BuffSize for enable buffer. set 0 to disable buffer
	BuffSize int `json:"buff_size" yaml:"buff_size"`
	// Levels for log record
	Levels []slog.Level `json:"levels" yaml:"levels"`
	// RotateTime for rotate file
	RotateTime rotatefile.RotateTime `json:"rotate_time" yaml:"rotate_time"`
	// MaxSize on rotate file by size.
	MaxSize uint64 `json:"max_size" yaml:"max_size"`
	// RenameFunc build filename for rotate file
	RenameFunc func(filepath string, rotateNum uint) string
}

// With more config settings func
func (c *Config) With(fns ...ConfigFn) *Config {
	for _, fn := range fns {
		fn(c)
	}
	return c
}

// SyncCloseWriter build by config
func (c *Config) SyncCloseWriter() (output SyncCloseWriter, err error) {
	output, err = fsutil.QuickOpenFile(c.Logfile)

	// wrap buffer
	if c.BuffSize > 0 {
		if c.BuffMode == BuffModeLine {
			output = bufwrite.NewLineWriterSize(output, c.BuffSize)
		} else {
			output = bufwrite.NewBufIOWriterSize(output, c.BuffSize)
		}
	}
	return
}

// RotateWriter build rotate writer by config
func (c *Config) RotateWriter() (output FlushCloseWriter, err error) {
	// create a rotate config.
	rc := rotatefile.NewConfig(c.Logfile)
	rc.MaxSize = c.MaxSize
	rc.CloseLock = true // has locked on logger.write()

	rc.RotateTime = c.RotateTime
	if c.RenameFunc != nil {
		rc.RenameFunc = c.RenameFunc
	}

	// create a rotating writer
	output, err = rc.Create()
	if err != nil {
		return nil, err
	}

	// wrap buffer
	if c.BuffSize > 0 {
		if c.BuffMode == BuffModeLine {
			output = bufwrite.NewLineWriterSize(output, c.BuffSize)
		} else {
			output = bufwrite.NewBufIOWriterSize(output, c.BuffSize)
		}
	}
	return
}

// NewEmptyConfig new config instance
func NewEmptyConfig(fns ...ConfigFn) *Config {
	c := &Config{
		Levels: slog.AllLevels,
	}
	return c.With(fns...)
}

// NewConfig new config instance with some default settings.
func NewConfig(fns ...ConfigFn) *Config {
	c := &Config{
		MaxSize:  rotatefile.DefaultMaxSize,
		BuffMode: BuffModeLine,
		BuffSize: DefaultBufferSize,
		Levels:   slog.AllLevels,
		// time rotate settings
		RotateTime: rotatefile.EveryHour,
	}

	return c.With(fns...)
}

// WithLogfile setting
func WithLogfile(logfile string) ConfigFn {
	return func(c *Config) {
		c.Logfile = logfile
	}
}

// WithRotateTime setting
func WithRotateTime(rt rotatefile.RotateTime) ConfigFn {
	return func(c *Config) {
		c.RotateTime = rt
	}
}

// WithBuffMode setting
func WithBuffMode(buffMode string) ConfigFn {
	return func(c *Config) {
		c.BuffMode = buffMode
	}
}

// WithBuffSize setting
func WithBuffSize(buffSize int) ConfigFn {
	return func(c *Config) {
		c.BuffSize = buffSize
	}
}

// WithMaxSize setting
func WithMaxSize(maxSize int) ConfigFn {
	return func(c *Config) {
		c.MaxSize = uint64(maxSize)
	}
}

// WithUseJSON setting
func WithUseJSON(useJson bool) ConfigFn {
	return func(c *Config) {
		c.UseJSON = useJson
	}
}

// WithLogLevels setting
func WithLogLevels(levels slog.Levels) ConfigFn {
	return func(c *Config) {
		c.Levels = levels
	}
}

// Builder struct for create handler
type Builder struct {
	Output   io.Writer
	Filepath string
	BuffSize int
	Levels   []slog.Level
}

// NewBuilder create
func NewBuilder() *Builder {
	return &Builder{}
}

// Build slog handler.
func (b *Builder) reset() {
	b.Output = nil
	b.Levels = b.Levels[:0]
	b.Filepath = ""
	b.BuffSize = 0
}

// Build slog handler.
func (b *Builder) Build() slog.Handler {
	defer b.reset()

	if b.Output != nil {
		return b.buildFromWriter(b.Output)
	}

	if b.Filepath != "" {
		f, err := QuickOpenFile(b.Filepath)
		if err != nil {
			panic(err)
		}

		return b.buildFromWriter(f)
	}

	panic("missing some information for build handler")
}

// Build slog handler.
func (b *Builder) buildFromWriter(w io.Writer) slog.Handler {
	if scw, ok := w.(SyncCloseWriter); ok {
		return NewSyncCloseHandler(scw, b.Levels)
	}

	if fcw, ok := w.(FlushCloseWriter); ok {
		return NewFlushCloseHandler(fcw, b.Levels)
	}

	if wc, ok := w.(io.WriteCloser); ok {
		return NewWriteCloser(wc, b.Levels)
	}

	return NewIOWriter(w, b.Levels)
}
