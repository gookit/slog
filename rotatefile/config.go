package rotatefile

import (
	"os"
	"strconv"
	"time"

	"github.com/gookit/goutil/timex"
)

type rotateLevel uint8

const (
	levelDay rotateLevel = iota
	levelHour
	levelMin
	levelSec
)

// RotateTime type
type RotateTime int

const (
	EveryDay    RotateTime = timex.OneDaySec
	EveryHour   RotateTime = timex.OneHourSec
	Every30Min  RotateTime = 30 * timex.OneMinSec
	Every15Min  RotateTime = 15 * timex.OneMinSec
	EveryMinute RotateTime = timex.OneMinSec
	EverySecond RotateTime = 1 // only use for tests
)

// Interval get check interval time
func (rt RotateTime) Interval() int64 {
	return int64(rt)
}

// level for rotate time
func (rt RotateTime) level() rotateLevel {
	switch {
	case rt > timex.OneDaySec:
		return levelDay
	case rt > timex.OneHourSec:
		return levelHour
	case rt > EveryMinute:
		return levelMin
	case rt > EverySecond:
		return levelSec
	default:
		return levelHour
	}
}

// TimeFormat get log file suffix format
func (rt RotateTime) TimeFormat() (suffixFormat string) {
	suffixFormat = "20060102_1500" // default is levelHour
	switch rt.level() {
	case levelDay:
		suffixFormat = "20060102"
	case levelHour:
		suffixFormat = "20060102_1500"
	case levelMin:
		suffixFormat = "20060102_1504"
	case levelSec:
		suffixFormat = "20060102_150405"
	}
	return
}

// String rotate type to string
func (rt RotateTime) String() string {
	switch rt.level() {
	case levelDay:
		return "Every Day"
	case levelMin:
		return "Every Minute"
	case levelSec:
		return "Every Second"
	case levelHour:
		return "Every Hours"
	default:
		return "Every Hours"
	}
}

// Config struct
type Config struct {
	// Filepath the log file path, will be rotating
	Filepath string `json:"filepath" yaml:"filepath"`
	// MaxSize file contents max size.
	//
	// unit is MB(megabytes)
	// default see DefaultMaxSize
	MaxSize uint64 `json:"max_size" yaml:"max_size"`
	// RotateTime the file rotate interval time.
	//
	// default is 1 hours. see EveryHour
	RotateTime RotateTime `json:"rotate_time" yaml:"rotate_time"`
	// CloseLock use sync lock on write contents, rotating file.
	//
	// default: false
	CloseLock bool `json:"close_lock" yaml:"close_lock"`
	// RenameFunc you can custom-build filename for rotate file by size.
	//
	// default see DefaultFilenameFunc
	RenameFunc func(filePath string, rotateNum uint) string
	// BackupNum max number for keep old files, 0 is not limit.
	BackupNum int `json:"backup_num" yaml:"backup_num"`
	// BackupTime max time for keep old files, 0 is not limit.
	//
	// unit is hours
	BackupTime int `json:"backup_time" yaml:"backup_time"`
}

// Create new RotateDispatcher
func (c *Config) Create() *RotateDispatcher {
	return New(c)
}

var (
	// DefaultMaxSize is the maximum size of a log file in bytes.
	//
	// unit is MB(megabytes)
	DefaultMaxSize uint64 = 1024 * 1024 * 1800
	// DefaultFilePerm perm and flags for create log file
	DefaultFilePerm  = 0664
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// DefaultFilenameFunc default new filename func
	DefaultFilenameFunc = func(filepath string, rotateNum uint) string {
		suffix := time.Now().Format("010215")

		// eg: /tmp/error.log => /tmp/error.log.163021_1
		return filepath + "." + suffix + "_" + strconv.Itoa(int(rotateNum))
	}
)

// NewDefaultConfig instance
func NewDefaultConfig() *Config {
	return &Config{
		MaxSize:    DefaultMaxSize,
		BackupNum:  20,
		RotateTime: EveryHour,
		RenameFunc: DefaultFilenameFunc,
	}
}

// NewConfig by file path
func NewConfig(filePath string) *Config {
	c := NewDefaultConfig()
	c.Filepath = filePath

	return c
}

// NewConfigWith custom func
func NewConfigWith(fn func(c *Config)) *Config {
	c := NewDefaultConfig()
	fn(c)
	return c
}
