// Package rotatefile is kept for backward compatibility only.
//
// Deprecated: rotatefile has moved to its own module. Use
// github.com/gookit/rotatefile instead. This package only re-exports the new
// location and will be removed in a future major version.
package rotatefile

import "github.com/gookit/rotatefile"

// ------- types (aliases keep type identity, methods/fields preserved) -------

type (
	// Deprecated: use github.com/gookit/rotatefile.Writer
	Writer = rotatefile.Writer
	// Deprecated: use github.com/gookit/rotatefile.Config
	Config = rotatefile.Config
	// Deprecated: use github.com/gookit/rotatefile.ConfigFn
	ConfigFn = rotatefile.ConfigFn
	// Deprecated: use github.com/gookit/rotatefile.RotateWriter
	RotateWriter = rotatefile.RotateWriter
	// Deprecated: use github.com/gookit/rotatefile.RotateTime
	RotateTime = rotatefile.RotateTime
	// Deprecated: use github.com/gookit/rotatefile.RotateMode
	RotateMode = rotatefile.RotateMode
	// Deprecated: use github.com/gookit/rotatefile.Clocker
	Clocker = rotatefile.Clocker
	// Deprecated: use github.com/gookit/rotatefile.ClockFn
	ClockFn = rotatefile.ClockFn
	// Deprecated: use github.com/gookit/rotatefile.MockClocker
	MockClocker = rotatefile.MockClocker
	// Deprecated: use github.com/gookit/rotatefile.FilesClear
	FilesClear = rotatefile.FilesClear
	// Deprecated: use github.com/gookit/rotatefile.CConfig
	CConfig = rotatefile.CConfig
	// Deprecated: use github.com/gookit/rotatefile.CConfigFunc
	CConfigFunc = rotatefile.CConfigFunc
)

// ------- constructors & functions -------

var (
	// Deprecated: use github.com/gookit/rotatefile.NewWriter
	NewWriter = rotatefile.NewWriter
	// Deprecated: use github.com/gookit/rotatefile.NewWriterWith
	NewWriterWith = rotatefile.NewWriterWith
	// Deprecated: use github.com/gookit/rotatefile.NewConfig
	NewConfig = rotatefile.NewConfig
	// Deprecated: use github.com/gookit/rotatefile.NewConfigWith
	NewConfigWith = rotatefile.NewConfigWith
	// Deprecated: use github.com/gookit/rotatefile.NewDefaultConfig
	NewDefaultConfig = rotatefile.NewDefaultConfig
	// Deprecated: use github.com/gookit/rotatefile.EmptyConfigWith
	EmptyConfigWith = rotatefile.EmptyConfigWith
	// Deprecated: use github.com/gookit/rotatefile.NewCConfig
	NewCConfig = rotatefile.NewCConfig
	// Deprecated: use github.com/gookit/rotatefile.NewFilesClear
	NewFilesClear = rotatefile.NewFilesClear
	// Deprecated: use github.com/gookit/rotatefile.NewMockClock
	NewMockClock = rotatefile.NewMockClock
	// Deprecated: use github.com/gookit/rotatefile.StringToRotateTime
	StringToRotateTime = rotatefile.StringToRotateTime
	// Deprecated: use github.com/gookit/rotatefile.StringToRotateMode
	StringToRotateMode = rotatefile.StringToRotateMode
	// Deprecated: use github.com/gookit/rotatefile.WithFilepath
	WithFilepath = rotatefile.WithFilepath
	// Deprecated: use github.com/gookit/rotatefile.WithCompress
	WithCompress = rotatefile.WithCompress
	// Deprecated: use github.com/gookit/rotatefile.WithDebugMode
	WithDebugMode = rotatefile.WithDebugMode
	// Deprecated: use github.com/gookit/rotatefile.WithBackupNum
	WithBackupNum = rotatefile.WithBackupNum
)

// ------- default vars -------

var (
	// Deprecated: use github.com/gookit/rotatefile.DefaultFilePerm
	DefaultFilePerm = rotatefile.DefaultFilePerm
	// Deprecated: use github.com/gookit/rotatefile.DefaultFileFlags
	DefaultFileFlags = rotatefile.DefaultFileFlags
	// Deprecated: use github.com/gookit/rotatefile.DefaultTimeClockFn
	DefaultTimeClockFn = rotatefile.DefaultTimeClockFn
)

// ------- constants -------

const (
	// Deprecated: use github.com/gookit/rotatefile.OneMByte
	OneMByte = rotatefile.OneMByte
	// Deprecated: use github.com/gookit/rotatefile.DefaultMaxSize
	DefaultMaxSize = rotatefile.DefaultMaxSize
	// Deprecated: use github.com/gookit/rotatefile.DefaultBackNum
	DefaultBackNum = rotatefile.DefaultBackNum
	// Deprecated: use github.com/gookit/rotatefile.DefaultBackTime
	DefaultBackTime = rotatefile.DefaultBackTime

	// Deprecated: use github.com/gookit/rotatefile.EveryMonth
	EveryMonth = rotatefile.EveryMonth
	// Deprecated: use github.com/gookit/rotatefile.EveryDay
	EveryDay = rotatefile.EveryDay
	// Deprecated: use github.com/gookit/rotatefile.EveryHour
	EveryHour = rotatefile.EveryHour
	// Deprecated: use github.com/gookit/rotatefile.Every30Min
	Every30Min = rotatefile.Every30Min
	// Deprecated: use github.com/gookit/rotatefile.Every15Min
	Every15Min = rotatefile.Every15Min
	// Deprecated: use github.com/gookit/rotatefile.EveryMinute
	EveryMinute = rotatefile.EveryMinute
	// Deprecated: use github.com/gookit/rotatefile.EverySecond
	EverySecond = rotatefile.EverySecond

	// Deprecated: use github.com/gookit/rotatefile.ModeRename
	ModeRename = rotatefile.ModeRename
	// Deprecated: use github.com/gookit/rotatefile.ModeCreate
	ModeCreate = rotatefile.ModeCreate
)
