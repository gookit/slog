package handler

import (
	"os"
	"sync"

	"github.com/gookit/goutil/fsutil"
)

// NopFlushClose no operation.
//
// provide empty Flush(), Close() methods, useful for tests.
type NopFlushClose struct{}

// Flush logs to disk
func (h *NopFlushClose) Flush() error {
	return nil
}

// Close handler
func (h *NopFlushClose) Close() error {
	return nil
}

type lockWrapper struct {
	sync.Mutex
	disable bool
}

// Lock it
func (lw *lockWrapper) Lock() {
	if false == lw.disable {
		lw.Mutex.Lock()
	}
}

// Unlock it
func (lw *lockWrapper) Unlock() {
	if !lw.disable {
		lw.Mutex.Unlock()
	}
}

// EnableLock enable lock
func (lw *lockWrapper) EnableLock(enable bool) {
	lw.disable = false == enable
}

// LockEnabled status
func (lw *lockWrapper) LockEnabled() bool {
	return lw.disable == false
}

// QuickOpenFile like os.OpenFile
func QuickOpenFile(filepath string) (*os.File, error) {
	return fsutil.OpenFile(filepath, DefaultFileFlags, DefaultFilePerm)
}
