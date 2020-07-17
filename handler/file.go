package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gookit/slog"
)

var onceLogDir sync.Once

var (
	// program pid
	pid = os.Getpid()
	// program name
	pName = filepath.Base(os.Args[0])
	host  = "unknownHost" // TODO
	// userName = "unknownuser"
)

// defaultMaxSize is the maximum size of a log file in bytes.
const defaultMaxSize uint64 = 1024 * 1024 * 1800

// FileHandler definition
type FileHandler struct {
	BaseHandler

	fpath string
	file  *os.File

	useJSON bool
	// perm for create log file
	FilePerm int
	// file contents max size
	MaxSize uint64
}

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(fpath string) *FileHandler {
	return NewFileHandler(fpath, true)
}

// NewFileHandler create new FileHandler
func NewFileHandler(fpath string, useJSON bool) *FileHandler {
	h := &FileHandler{
		fpath:   fpath,
		useJSON: useJSON,
		MaxSize: defaultMaxSize,
	}

	if useJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	}

	return h
}

// Writer return *os.File
func (h *FileHandler) Writer() io.Writer {
	return h.file
}

// Sync logs to disk file
func (h *FileHandler) Sync() error {
	return h.file.Sync()
}

// Handle the log record
func (h *FileHandler) Handle(r *slog.Record) error {
	bts, err := h.Formatter().Format(r)
	if err != nil {
		return err
	}

	_, err = h.file.Write(bts)
	return err
}

func logName(tag string, t time.Time) string {
	// return tag + t.Nanosecond()
	return tag + t.Format("20060102T15:04:05Z07:00")
}

// from glog
// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}

	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	// TODO ...
	// onceLogDir.Do(func() {
	// 	fsutil.Mkdir(filepath.Dir(h.fpath))
	// })

	dir := "some"

	// name, link := logName(tag, t)
	name := "xxx.log"
	var lastErr error

	fName := filepath.Join(dir, name)
	f, err = os.Create(fName)
	if err == nil {
		// symlink := filepath.Join(dir, link)
		// os.Remove(symlink)        // ignore err
		// os.Symlink(name, symlink) // ignore err
		return f, fName, nil
	}
	lastErr = err

	return nil, "", fmt.Errorf("log: cannot create log file, error: %v", lastErr)
}
