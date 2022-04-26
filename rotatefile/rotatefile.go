package rotatefile

import (
	"io"
	"sync"
)

// RotateFiler interface
type RotateFiler interface {
	Flush() error
	Sync() error
	// WriteCloser the output writer
	io.WriteCloser
}

// RotateFiles multi files. TODO
// use for rotate and clear other program produce log files
//
// refer file-rotatelogs
type RotateFiles struct {
	sync.Mutex
	cfg     *Config
	pattern string
}
