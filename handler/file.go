package handler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// var onceLogDir sync.Once

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
	*bufio.Writer // has Flash() method
	filepath      string
	file          *os.File

	MaxSize uint64
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		MaxSize: defaultMaxSize,
	}
}

func (h *FileHandler) Flush() error {
	return h.Writer.Flush()
}

func (h *FileHandler) Sync() error {
	return h.file.Sync()
}

func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	// TODO ...
	// onceLogDir.Do(createLogDir)

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
