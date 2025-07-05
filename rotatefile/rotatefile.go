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

const (
	// OneMByte size
	OneMByte uint64 = 1024 * 1024

	// DefaultMaxSize of a log file. default is 20M.
	DefaultMaxSize = 20 * OneMByte
	// DefaultBackNum default backup numbers for old files.
	DefaultBackNum uint = 20
	// DefaultBackTime default backup time for old files. default keeps a week.
	DefaultBackTime uint = 24 * 7
)
