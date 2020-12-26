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
	host  = "unknownHost" // TODO
	// userName = "unknownuser"
)

// defaultMaxSize is the maximum size of a log file in bytes.
const defaultMaxSize uint64 = 1024 * 1024 * 1800

const (
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

// FileHandler definition
type FileHandler struct {
	LevelsWithFormatter

	mu sync.Mutex

	// log file path. eg: "/var/log/my-app.log"
	fpath string
	file  *os.File
	bufio *bufio.Writer

	useJSON bool
	// NoBuffer on write log records
	NoBuffer bool
	// FileFlag for create. default: os.O_CREATE|os.O_WRONLY|os.O_APPEND
	FileFlag int
	// FileMode perm for create log file. (it's os.FileMode)
	FileMode uint32
	// BuffSize for enable buffer
	BuffSize int
	// file contents max size
	MaxSize uint64
	// RenameFunc for rotate file
	RenameFunc func(fpath string) string
}

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(fpath string) *FileHandler {
	return NewFileHandler(fpath, true)
}

// NewFileHandler create new FileHandler
func NewFileHandler(fpath string, useJSON bool) *FileHandler {
	h := &FileHandler{
		fpath:    fpath,
		useJSON:  useJSON,
		MaxSize:  defaultMaxSize,
		FileMode: 0664, // default FileMode
		FileFlag: os.O_CREATE | os.O_WRONLY | os.O_APPEND,
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

	return h
}

// Configure the handler
func (h *FileHandler) Configure(fn func(h *FileHandler)) *FileHandler {
	fn(h)
	return h
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

	h.mu.Lock()
	defer h.mu.Unlock()

	bts, err = h.Formatter().Format(r)
	if err != nil {
		return
	}

	// create file
	if h.file == nil {
		dPath := path.Dir(h.fpath)
		err = os.MkdirAll(dPath, 0777)
		if err != nil {
			return
		}

		h.file, err = os.OpenFile(h.fpath, h.FileFlag, os.FileMode(h.FileMode))
		if err != nil {
			return
		}
	}

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
