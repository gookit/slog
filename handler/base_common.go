package handler

import (
	"io"
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

type fileWrapper struct {
	path string
	file *os.File
}

// ReopenFile the log file
func (h *fileWrapper) ReopenFile() error {
	if h.file != nil {
		h.file.Close()
	}

	file, err := fsutil.QuickOpenFile(h.path)
	if err != nil {
		return err
	}

	h.file = file
	return err
}

// Write contents to *os.File
func (h *fileWrapper) Write(bts []byte) (n int, err error) {
	return h.file.Write(bts)
}

// Writer return *os.File
func (h *fileWrapper) Writer() io.Writer {
	return h.file
}

// Close handler, will be flush logs to file, then close file
func (h *fileWrapper) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}
	return h.file.Close()
}

// Flush logs to disk file
func (h *fileWrapper) Flush() error {
	return h.file.Sync()
}
