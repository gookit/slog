package rotatefile

import (
	"fmt"
	"os"
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

// RotateTime for rotate file. unit is seconds.
//
// EveryDay:
//   - "error.log.20201223"
//
// EveryHour, Every30Min, EveryMinute:
//   - "error.log.20201223_1500"
//   - "error.log.20201223_1530"
//   - "error.log.20201223_1523"
type RotateTime int

// built in rotate time constants
const (
	EveryMonth  RotateTime = 30 * timex.OneDaySec
	EveryDay    RotateTime = timex.OneDaySec
	EveryHour   RotateTime = timex.OneHourSec
	Every30Min  RotateTime = 30 * timex.OneMinSec
	Every15Min  RotateTime = 15 * timex.OneMinSec
	EveryMinute RotateTime = timex.OneMinSec
	EverySecond RotateTime = 1 // only use for tests
)

// Interval get check interval time. unit is seconds.
func (rt RotateTime) Interval() int64 {
	return int64(rt)
}

// FirstCheckTime for rotate file.
// - will automatically align the time from the start of each hour.
func (rt RotateTime) FirstCheckTime(now time.Time) time.Time {
	interval := rt.Interval()

	switch rt.level() {
	case levelDay:
		return timex.DayEnd(now)
	case levelHour:
		// should check on H:59:59.500
		return timex.HourStart(now).Add(timex.OneHour - 500*time.Millisecond)
	case levelMin:
		// eg: minutes=5
		minutes := int(interval / 60)
		nextMin := now.Minute() + minutes

		// will rotate at next hour start. eg: now.Minute()=57, nextMin=62.
		if nextMin >= 60 {
			return timex.HourStart(now).Add(timex.OneHour)
		}

		// eg: now.Minute()=37, nextMin=42, will get nextDur=40
		nextDur := time.Duration(nextMin).Round(time.Duration(minutes))
		return timex.HourStart(now).Add(nextDur * time.Minute)
	default: // levelSec
		return now.Add(time.Duration(interval) * time.Second)
	}
}

// level for rotate time
func (rt RotateTime) level() rotateLevel {
	switch {
	case rt >= timex.OneDaySec:
		return levelDay
	case rt >= timex.OneHourSec:
		return levelHour
	case rt >= EveryMinute:
		return levelMin
	default:
		return levelSec
	}
}

// TimeFormat get log file suffix format
//
// EveryDay:
//   - "error.log.20201223"
//
// EveryHour, Every30Min, EveryMinute:
//   - "error.log.20201223_1500"
//   - "error.log.20201223_1530"
//   - "error.log.20201223_1523"
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
		return fmt.Sprintf("Every %d Day", rt.Interval()/timex.OneDaySec)
	case levelHour:
		return fmt.Sprintf("Every %d Hours", rt.Interval()/timex.OneHourSec)
	case levelMin:
		return fmt.Sprintf("Every %d Minutes", rt.Interval()/timex.OneMinSec)
	default: // levelSec
		return fmt.Sprintf("Every %d Seconds", rt.Interval())
	}
}

// Clocker is the interface used for determine the current time
type Clocker interface {
	Now() time.Time
}

// ClockFn func
type ClockFn func() time.Time

// Now implements the Clocker
func (fn ClockFn) Now() time.Time {
	return fn()
}

// ConfigFn for setting config
type ConfigFn func(c *Config)

// Config struct for rotate dispatcher
type Config struct {
	// Filepath the log file path, will be rotating. eg: "logs/error.log"
	Filepath string `json:"filepath" yaml:"filepath"`

	// FilePerm for create log file. default DefaultFilePerm
	FilePerm os.FileMode `json:"file_perm" yaml:"file_perm"`

	// RotateMode for rotate file. default ModeRename
	RotateMode RotateMode `json:"rotate_mode" yaml:"rotate_mode"`

	// MaxSize file contents max size, unit is bytes.
	// If is equals zero, disable rotate file by size
	//
	// default see DefaultMaxSize
	MaxSize uint64 `json:"max_size" yaml:"max_size"`

	// RotateTime the file rotate interval time, unit is seconds.
	// If is equals zero, disable rotate file by time
	//
	// default: EveryHour
	RotateTime RotateTime `json:"rotate_time" yaml:"rotate_time"`

	// CloseLock use sync lock on write contents, rotating file.
	//
	// default: false
	CloseLock bool `json:"close_lock" yaml:"close_lock"`

	// BackupNum max number for keep old files.
	//
	// 0 is not limit, default is DefaultBackNum
	BackupNum uint `json:"backup_num" yaml:"backup_num"`

	// BackupTime max time for keep old files, unit is hours.
	//
	// 0 is not limit, default is DefaultBackTime
	BackupTime uint `json:"backup_time" yaml:"backup_time"`

	// Compress determines if the rotated log files should be compressed using gzip.
	// The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`

	// RenameFunc you can custom-build filename for rotate file by size.
	//
	// Example:
	//
	//  c.RenameFunc = func(filepath string, rotateNum uint) string {
	// 		suffix := time.Now().Format("06010215")
	//
	// 		// eg: /tmp/error.log => /tmp/error.log.24032116_894136
	// 		return filepath + fmt.Sprintf(".%s_%d", suffix, rotateNum)
	//  }
	RenameFunc func(filePath string, rotateNum uint) string

	// TimeClock for rotate file by time.
	TimeClock Clocker

	// DebugMode for debug on development.
	DebugMode bool
}

func (c *Config) backupDuration() time.Duration {
	if c.BackupTime < 1 {
		return 0
	}
	return time.Duration(c.BackupTime) * time.Hour
}

// With more config setting func
func (c *Config) With(fns ...ConfigFn) *Config {
	for _, fn := range fns {
		fn(c)
	}
	return c
}

// Create new Writer by config
func (c *Config) Create() (*Writer, error) { return NewWriter(c) }

// IsMode check rotate mode
func (c *Config) IsMode(m RotateMode) bool { return c.RotateMode == m }

var (
	// DefaultFilePerm perm and flags for create log file
	DefaultFilePerm os.FileMode = 0664
	// DefaultFileFlags for open log file
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// DefaultTimeClockFn for create time
	DefaultTimeClockFn = ClockFn(func() time.Time {
		return time.Now()
	})
)

// NewDefaultConfig instance
func NewDefaultConfig() *Config {
	return &Config{
		MaxSize:    DefaultMaxSize,
		RotateTime: EveryHour,
		BackupNum:  DefaultBackNum,
		BackupTime: DefaultBackTime,
		// RenameFunc: DefaultFilenameFn,
		TimeClock: DefaultTimeClockFn,
		FilePerm:  DefaultFilePerm,
	}
}

// NewConfig by file path, and can with custom setting
func NewConfig(filePath string, fns ...ConfigFn) *Config {
	if len(fns) == 0 {
		return NewConfigWith(WithFilepath(filePath))
	}
	return NewConfigWith(append(fns, WithFilepath(filePath))...)
}

// NewConfigWith custom func
func NewConfigWith(fns ...ConfigFn) *Config {
	return NewDefaultConfig().With(fns...)
}

// EmptyConfigWith new empty config with custom func
func EmptyConfigWith(fns ...ConfigFn) *Config {
	c := &Config{
		// RenameFunc: DefaultFilenameFn,
		TimeClock: DefaultTimeClockFn,
		FilePerm:  DefaultFilePerm,
	}

	return c.With(fns...)
}

// WithFilepath setting
func WithFilepath(logfile string) ConfigFn {
	return func(c *Config) {
		c.Filepath = logfile
	}
}

// WithDebugMode setting for debug mode
func WithDebugMode(c *Config) {
	c.DebugMode = true
}

// WithCompress setting for compress
func WithCompress(c *Config) {
	c.Compress = true
}

// WithBackupNum setting for backup number
func WithBackupNum(num uint) ConfigFn {
	return func(c *Config) {
		c.BackupNum = num
	}
}
