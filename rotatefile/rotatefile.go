// Package rotatefile provides simple file rotation, compression and cleanup.
package rotatefile

import (
	"io"
)

// RotateWriter interface
type RotateWriter interface {
	io.WriteCloser
	Clean() error
	Flush() error
	Rotate() error
	Sync() error
}

// RotateMode for rotate file. 0: rename, 1: create
type RotateMode uint8

// String get string name
func (m RotateMode) String() string {
	switch m {
	case ModeRename:
		return "rename"
	case ModeCreate:
		return "create"
	default:
		return "unknown"
	}
}

const (
	// ModeRename rotating file by rename.
	//
	// Example flow:
	//  - always write to "error.log"
	//  - rotating by rename it to "error.log.20201223"
	//  - then re-create "error.log"
	ModeRename RotateMode = iota

	// ModeCreate rotating file by create new file.
	//
	// Example flow:
	//  - directly create new file on each rotate time. eg: "error.log.20201223", "error.log.20201224"
	ModeCreate
)

const (
	// OneMByte size
	OneMByte uint64 = 1024 * 1024

	// DefaultMaxSize of a log file. default is 20M.
	DefaultMaxSize = 20 * OneMByte
	// DefaultBackNum default backup numbers for old files.
	DefaultBackNum uint = 20
	// DefaultBackTime default backup time for old files. default keep a week.
	DefaultBackTime uint = 24 * 7
)
