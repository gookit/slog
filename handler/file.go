package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
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

// MaxSize is the maximum size of a log file in bytes.
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

func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	// TODO ...
	onceLogDir.Do(func() {
		fsutil.Mkdir(filepath.Dir(h.fpath))
	})

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
