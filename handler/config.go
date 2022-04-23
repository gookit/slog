package handler

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/bufwrite"
	"github.com/gookit/slog/rotatefile"
)

// SimpleConfig struct
type SimpleConfig struct {
	Logfile string `json:"logfile" yaml:"logfile"`
	UseJSON bool   `json:"use_json"`
	// BuffSize for enable buffer
	BuffSize int `json:"buff_size"`
	// Levels for logger
	Levels []slog.Level `json:"levels" yaml:"levels"`
}

// NewSimpleConfig create
func NewSimpleConfig(fns ...func(c *SimpleConfig)) *SimpleConfig {
	c := &SimpleConfig{
		BuffSize: DefaultBufferSize,
		Levels:   slog.AllLevels,
	}

	if len(fns) > 0 {
		for _, fn := range fns {
			fn(c)
		}
	}
	return c
}

// Config struct
type Config struct {
	// Logfile for write logs
	Logfile string `json:"logfile" yaml:"logfile"`
	// RotateTime for rotate file
	RotateTime RotateTime `json:"rotate_time" yaml:"rotate_time"`
	// NoBuffer on write log records
	MaxSize uint `json:"no_buffer" yaml:"max_size"`
	// BuffSize for enable buffer
	BuffSize int `json:"buff_size" yaml:"buff_size"`
	// Levels for logger
	Levels []slog.Level `json:"levels" yaml:"levels"`
	// RenameFunc build filename for rotate file
	RenameFunc func(filepath string, rotateNum uint) string
}

// CreateWriter build writer by config
func (c *Config) CreateWriter() (output FlushCloseWriter, err error) {
	output, err = c.RotateConfig().Create()
	if err != nil {
		return nil, err
	}

	// wrap buffer writer
	if c.BuffSize > 0 {
		output = bufwrite.NewBufIOWriterSize(output, c.BuffSize)
	}
	return
}

// RotateConfig build
func (c *Config) RotateConfig() *rotatefile.Config {
	rc := rotatefile.NewConfig(c.Logfile)
	rc.CloseLock = true // lock is opened on logger.write()
	rc.RotateTime = c.RotateTime

	rc.MaxSize = c.MaxSize
	if c.RenameFunc != nil {
		rc.RenameFunc = c.RenameFunc
	}

	return rc
}

// NewConfig new config instance
func NewConfig(fns ...func(c *Config)) *Config {
	c := &Config{
		MaxSize:    rotatefile.DefaultMaxSize,
		BuffSize:   DefaultBufferSize,
		RotateTime: rotatefile.EveryHour,
	}

	for _, fn := range fns {
		fn(c)
	}
	return c
}
