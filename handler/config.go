package handler

import (
	"io"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
	"github.com/gookit/slog/rotatefile"
)

// the buff mode consts
const (
	BuffModeLine = "line"
	BuffModeBite = "bite"
)

const (
	// LevelModeList use level list for filter record write
	LevelModeList uint8 = iota
	// LevelModeValue use level value compare for filter record write
	LevelModeValue
)

// ConfigFn for config some settings
type ConfigFn func(c *Config)

// Config struct
type Config struct {
	// Logfile for write logs
	Logfile string `json:"logfile" yaml:"logfile"`

	// LevelMode for filter log record. default LevelModeList
	LevelMode uint8 `json:"level_mode" yaml:"level_mode"`

	// Level value. use on LevelMode = LevelModeValue
	Level slog.Level `json:"level" yaml:"level"`

	// Levels list for write. use on LevelMode = LevelModeList
	Levels []slog.Level `json:"levels" yaml:"levels"`

	// UseJSON for format logs
	UseJSON bool `json:"use_json" yaml:"use_json"`

	// BuffMode type name. allow: line, bite
	BuffMode string `json:"buff_mode" yaml:"buff_mode"`

	// BuffSize for enable buffer, unit is bytes. set 0 to disable buffer
	BuffSize int `json:"buff_size" yaml:"buff_size"`

	// RotateTime for rotate file, unit is seconds.
	RotateTime rotatefile.RotateTime `json:"rotate_time" yaml:"rotate_time"`

	// MaxSize on rotate file by size, unit is bytes.
	MaxSize uint64 `json:"max_size" yaml:"max_size"`

	// Compress determines if the rotated log files should be compressed using gzip.
	// The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`

	// BackupNum max number for keep old files.
	//
	// 0 is not limit, default is 20.
	BackupNum uint `json:"backup_num" yaml:"backup_num"`

	// BackupTime max time for keep old files, unit is hours.
	//
	// 0 is not limit, default is a week.
	BackupTime uint `json:"backup_time" yaml:"backup_time"`

	// RenameFunc build filename for rotate file
	RenameFunc func(filepath string, rotateNum uint) string
}

// NewEmptyConfig new config instance
func NewEmptyConfig(fns ...ConfigFn) *Config {
	c := &Config{Levels: slog.AllLevels}
	return c.WithConfigFn(fns...)
}

// NewConfig new config instance with some default settings.
func NewConfig(fns ...ConfigFn) *Config {
	c := &Config{
		Levels:   slog.AllLevels,
		BuffMode: BuffModeLine,
		BuffSize: DefaultBufferSize,
		// rotate file settings
		MaxSize:    rotatefile.DefaultMaxSize,
		RotateTime: rotatefile.EveryHour,
		// old files clean settings
		BackupNum:  rotatefile.DefaultBackNum,
		BackupTime: rotatefile.DefaultBackTime,
	}

	return c.WithConfigFn(fns...)
}

// With more config settings func
func (c *Config) With(fns ...ConfigFn) *Config {
	return c.WithConfigFn(fns...)
}

// WithConfigFn more config settings func
func (c *Config) WithConfigFn(fns ...ConfigFn) *Config {
	for _, fn := range fns {
		fn(c)
	}
	return c
}

func (c *Config) newLevelFormattable() slog.LevelFormattable {
	if c.LevelMode == LevelModeValue {
		return slog.NewLvFormatter(c.Level)
	}
	return slog.NewLvsFormatter(c.Levels)
}

// CreateHandler quick create a handler by config
func (c *Config) CreateHandler() (*SyncCloseHandler, error) {
	output, err := c.CreateWriter()
	if err != nil {
		return nil, err
	}

	h := &SyncCloseHandler{
		Output: output,
		// with log level and formatter
		LevelFormattable: c.newLevelFormattable(),
	}

	if c.UseJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	}
	return h, nil
}

// RotateWriter build rotate writer by config
func (c *Config) RotateWriter() (output SyncCloseWriter, err error) {
	if c.MaxSize == 0 && c.RotateTime == 0 {
		return nil, errorx.Raw("slog: cannot create rotate writer, MaxSize and RotateTime both is 0")
	}

	return c.CreateWriter()
}

// CreateWriter build writer by config
func (c *Config) CreateWriter() (output SyncCloseWriter, err error) {
	if c.Logfile == "" {
		return nil, errorx.Raw("slog: logfile cannot be empty for create writer")
	}

	// create a rotate config.
	if c.MaxSize > 0 || c.RotateTime > 0 {
		rc := rotatefile.EmptyConfigWith()

		// has locked on logger.write()
		rc.CloseLock = true
		rc.Filepath = c.Logfile
		// copy settings
		rc.MaxSize = c.MaxSize
		rc.RotateTime = c.RotateTime
		rc.BackupNum = c.BackupNum
		rc.BackupTime = c.BackupTime
		rc.Compress = c.Compress

		if c.RenameFunc != nil {
			rc.RenameFunc = c.RenameFunc
		}

		// create a rotating writer
		output, err = rc.Create()
	} else {
		output, err = fsutil.QuickOpenFile(c.Logfile)
	}

	if err != nil {
		return nil, err
	}

	// wrap buffer
	if c.BuffSize > 0 {
		output = c.wrapBuffer(output)
	}
	return
}

type flushSyncCloseWriter interface {
	FlushCloseWriter
	Sync() error
}

// wrap buffer for the writer
func (c *Config) wrapBuffer(w io.Writer) (bw flushSyncCloseWriter) {
	if c.BuffSize == 0 {
		panic("slog: buff size cannot be zero on wrap buffer")
	}

	if c.BuffMode == BuffModeLine {
		bw = bufwrite.NewLineWriterSize(w, c.BuffSize)
	} else {
		bw = bufwrite.NewBufIOWriterSize(w, c.BuffSize)
	}
	return bw
}

//
// ---------------------------------------------------------------------------
// global config func
// ---------------------------------------------------------------------------
//

// WithLogfile setting
func WithLogfile(logfile string) ConfigFn {
	return func(c *Config) { c.Logfile = logfile }
}

// WithLevelMode setting
func WithLevelMode(mode uint8) ConfigFn {
	return func(c *Config) { c.LevelMode = mode }
}

// WithLogLevel setting
func WithLogLevel(level slog.Level) ConfigFn {
	return func(c *Config) { c.Level = level }
}

// WithLogLevels setting
func WithLogLevels(levels slog.Levels) ConfigFn {
	return func(c *Config) { c.Levels = levels }
}

// WithLevelNames set levels by level names.
func WithLevelNames(names []string) ConfigFn {
	levels := make([]slog.Level, 0, len(names))
	for _, name := range names {
		levels = append(levels, slog.LevelByName(name))
	}

	return func(c *Config) {
		c.Levels = levels
	}
}

// WithRotateTime setting
func WithRotateTime(rt rotatefile.RotateTime) ConfigFn {
	return func(c *Config) { c.RotateTime = rt }
}

// WithBackupNum setting
func WithBackupNum(n uint) ConfigFn {
	return func(c *Config) { c.BackupNum = n }
}

// WithBackupTime setting
func WithBackupTime(bt uint) ConfigFn {
	return func(c *Config) { c.BackupTime = bt }
}

// WithBuffMode setting
func WithBuffMode(buffMode string) ConfigFn {
	return func(c *Config) { c.BuffMode = buffMode }
}

// WithBuffSize setting
func WithBuffSize(buffSize int) ConfigFn {
	return func(c *Config) { c.BuffSize = buffSize }
}

// WithMaxSize setting
func WithMaxSize(maxSize uint64) ConfigFn {
	return func(c *Config) { c.MaxSize = maxSize }
}

// WithCompress setting
func WithCompress(compress bool) ConfigFn {
	return func(c *Config) { c.Compress = compress }
}

// WithUseJSON setting
func WithUseJSON(useJSON bool) ConfigFn {
	return func(c *Config) { c.UseJSON = useJSON }
}
