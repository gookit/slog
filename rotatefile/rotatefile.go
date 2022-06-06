package rotatefile

import (
	"fmt"
	"io"
	"os"
)

// RotateWriter interface
type RotateWriter interface {
	io.WriteCloser
	Clean() error
	Flush() error
	Rotate() error
	Sync() error
}

func printErrln(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}
