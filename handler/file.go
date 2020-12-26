package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
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
	hName = "unknownHost" // TODO
	// uName = "unknownUser"
)

var (
	// DefaultMaxSize is the maximum size of a log file in bytes.
	DefaultMaxSize uint64 = 1024 * 1024 * 1800
	// perm and flags for create log file
	DefaultFilePerm  = 0664
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

// FileHandler definition
type FileHandler struct {
	// fileWrapper
	lockWrapper
	// LevelsWithFormatter support limit log levels and formatter
	LevelsWithFormatter

	// log file path. eg: "/var/log/my-app.log"
	fpath string
	file  *os.File
	bufio *bufio.Writer

	useJSON bool
	// NoBuffer on write log records
	NoBuffer bool
	// BuffSize for enable buffer
	BuffSize int
}

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(filepath string) (*FileHandler, error) {
	return NewFileHandler(filepath, true)
}

// MustFileHandler create file handler
func MustFileHandler(filepath string, useJSON bool) *FileHandler {
	h, err := NewFileHandler(filepath, useJSON)
	if err != nil {
		panic(err)
	}

	return h
}

// NewFileHandler create new FileHandler
func NewFileHandler(filepath string, useJSON bool) (*FileHandler, error) {
	h := &FileHandler{
		fpath:   filepath,
		useJSON: useJSON,
		// MaxSize:  DefaultMaxSize,
		// FileMode: DefaultFilePerm, // default FileMode
		// FileFlag: DefaultFileFlags,
		// init log levels
		LevelsWithFormatter: LevelsWithFormatter{
			Levels: slog.AllLevels, // default log all levels
		},
	}

	if useJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	} else {
		h.SetFormatter(slog.NewTextFormatter())
	}

	file, err := openFile(filepath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return nil, err
	}

	h.file = file
	return h, nil
}

// Configure the handler
func (h *FileHandler) Configure(fn func(h *FileHandler)) *FileHandler {
	fn(h)
	return h
}

// ReopenFile the log file
func (h *FileHandler) ReopenFile() error {
	file, err := openFile(h.fpath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return err
	}

	h.file = file
	return err
}

// Writer return *os.File
func (h *FileHandler) Writer() io.Writer {
	return h.file
}

// Close handler, will be flush logs to file, then close file
func (h *FileHandler) Close() error {
	if err := h.Flush(); err != nil {
		return err
	}

	return h.file.Close()
}

// Flush logs to disk file
func (h *FileHandler) Flush() error {
	// flush buffers to h.file
	if h.bufio != nil {
		err := h.bufio.Flush()
		if err != nil {
			return err
		}
	}

	return h.file.Sync()
}

// Handle the log record
func (h *FileHandler) Handle(r *slog.Record) (err error) {
	var bts []byte
	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// if enable lock
	if h.LockEnabled() {
		h.Lock()
		defer h.Unlock()
	}

	// create file
	// if h.file == nil {
	// 	h.file, err = openFile(h.fpath, h.FileFlag, h.FileMode)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	// direct write logs
	if h.NoBuffer {
		_, err = h.file.Write(bts)
		return
	}

	// enable buffer
	if h.bufio == nil {
		h.bufio = bufio.NewWriterSize(h.file, h.BuffSize)
	}

	_, err = h.bufio.Write(bts)
	return
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
	// 	fsutil.Mkdir(fpath.Dir(h.fpath))
	// })

	dir := "some"

	// name, link := logName(tag, t)
	name := "xxx.log"
	var lastErr error

	fName := filepath.Join(dir, name)
	f, err = os.Create(fName)
	if err == nil {
		// symlink := fpath.Join(dir, link)
		// os.Remove(symlink)        // ignore err
		// os.Symlink(name, symlink) // ignore err
		return f, fName, nil
	}
	lastErr = err

	return nil, "", fmt.Errorf("log: cannot create log file, error: %v", lastErr)
}

func openFile(filepath string, flag int, mode int) (*os.File, error) {
	fileDir := path.Dir(filepath)

	// if err := os.Mkdir(dir, 0777); err != nil {
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, flag, os.FileMode(mode))
	if err != nil {
		return nil, err
	}

	return file, nil
}
