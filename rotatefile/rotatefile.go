package rotatefile

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// RotateWriter interface
type RotateWriter interface {
	io.WriteCloser
	Clean() error
	Flush() error
	Rotate() error
	Sync() error
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

func printlnStderr(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}
