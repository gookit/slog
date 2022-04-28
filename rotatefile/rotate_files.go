package rotatefile

import "sync"

// RotateFiles multi files. TODO
// use for rotate and clear other program produce log files
//
// refer file-rotatelogs
type RotateFiles struct {
	sync.Mutex
	cfg     *Config
	pattern string
}

// Rotate do rotate handle
func (r *RotateFiles) Rotate() error {
	return nil
}
