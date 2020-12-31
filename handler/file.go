package handler

import (
	"bufio"
	"io"
	"os"
	"path"

	"github.com/gookit/slog"
)

// TODO use this ... ?
// var onceLogDir sync.Once

// FileHandler definition
//
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

// func WithBuffer()  {
//
// }

// JSONFileHandler create new FileHandler with JSON formatter
func JSONFileHandler(logfile string) (*FileHandler, error) {
	return NewFileHandler(logfile, true)
}

// MustFileHandler create file handler
func MustFileHandler(logfile string, useJSON bool) *FileHandler {
	h, err := NewFileHandler(logfile, useJSON)
	if err != nil {
		panic(err)
	}

	return h
}

// NewFileHandler create new FileHandler
func NewFileHandler(logfile string, useJSON bool) (*FileHandler, error) {
	h := &FileHandler{
		fpath:    logfile,
		useJSON:  useJSON,
		BuffSize: defaultBufferSize,
		// FileMode: DefaultFilePerm, // default FileMode
		// FileFlag: DefaultFileFlags,

		// default log all levels
		LevelsWithFormatter: newLvsFormatter(slog.AllLevels),
	}

	if useJSON {
		h.SetFormatter(slog.NewJSONFormatter())
	} else {
		h.SetFormatter(slog.NewTextFormatter())
	}

	file, err := QuickOpenFile(h.fpath)
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
	file, err := OpenFile(h.fpath, DefaultFileFlags, DefaultFilePerm)
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
	h.Lock()
	defer h.Unlock()

	// create file
	// if h.file == nil {
	// 	h.file, err = OpenFile(h.fpath, h.FileFlag, h.FileMode)
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
	if h.bufio == nil && h.BuffSize > 0 {
		h.bufio = bufio.NewWriterSize(h.file, h.BuffSize)
	}

	_, err = h.bufio.Write(bts)
	return
}

// QuickOpenFile like os.OpenFile
func QuickOpenFile(filepath string) (*os.File, error) {
	return OpenFile(filepath, DefaultFileFlags, DefaultFilePerm)
}

// OpenFile like os.OpenFile
func OpenFile(filepath string, flag int, perm int) (*os.File, error) {
	fileDir := path.Dir(filepath)

	// if err := os.Mkdir(dir, 0777); err != nil {
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, flag, os.FileMode(perm))
	if err != nil {
		return nil, err
	}

	return file, nil
}
