package rotatefile

import (
	"os"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// RotateDispatcher struct
type RotateDispatcher struct {
	*Config
	sync.Mutex
	file *os.File

	// context use for rotating file by size
	written   uint64 // written size
	rotateNum uint   // rotate times number

	// context use for rotating file by time
	suffixFormat   string // the rotating file name suffix.
	checkInterval  int64
	nextRotatingAt int64
}

// New create rotate dispatcher with config.
func New(c *Config) *RotateDispatcher {
	d := &RotateDispatcher{
		Config: c,
	}

	if err := d.Init(); err != nil {
		panic(err)
	}
	return d
}

// Init rotate dispatcher
func (d *RotateDispatcher) Init() error {
	d.suffixFormat = d.RotateTime.TimeFormat()
	d.checkInterval = d.RotateTime.Interval()

	// open the log file
	file, err := fsutil.QuickOpenFile(d.Filepath)
	if err != nil {
		return err
	}

	// storage next rotating time
	fStat, err := file.Stat()
	if err != nil {
		return err
	}

	d.file = file
	d.nextRotatingAt = fStat.ModTime().Unix() + d.checkInterval
	return nil
}

// ReopenFile the log file
func (d *RotateDispatcher) ReopenFile() error {
	if d.file != nil {
		d.file.Close()
	}

	file, err := fsutil.QuickOpenFile(d.Filepath)
	if err != nil {
		return err
	}

	d.file = file
	return nil
}

// Rotate file handle
func (d *RotateDispatcher) Rotate() error {
	return nil
}

func (d *RotateDispatcher) Flush() error {
	return nil
}

func (d *RotateDispatcher) Sync() error {
	return nil
}

func (d *RotateDispatcher) Close() error {
	return nil
}

func (d *RotateDispatcher) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

func (d *RotateDispatcher) Write(p []byte) (n int, err error) {
	// if enable lock
	if d.CloseLock == false {
		d.Lock()
		defer d.Unlock()
	}

	n, err = d.file.Write(p)
	if err != nil {
		return
	}

	d.written += uint64(n)

	// do rotate file by time
	if d.checkInterval > 0 {
		err = d.rotatingByTime()
		if err != nil {
			return
		}
	}

	// do rotate file by size
	if d.MaxSize > 0 && d.written >= d.MaxSize {
		err = d.rotatingBySize()
	}
	return
}

func (d *RotateDispatcher) rotatingByTime() error {
	now := time.Now()
	if d.nextRotatingAt > now.Unix() {
		return nil
	}

	// rename current to new file.
	// eg: /tmp/error.log => /tmp/error.log.20220423_1600
	newFilepath := d.Filepath + "." + now.Format(d.suffixFormat)

	err := d.rotatingFile(newFilepath)

	// storage next rotating time
	d.nextRotatingAt = now.Unix() + d.checkInterval
	return err
}

func (d *RotateDispatcher) rotatingBySize() error {
	// rename current to new file
	d.rotateNum++

	// eg: /tmp/error.log => /tmp/error.log.163021_1
	newFilepath := d.RenameFunc(d.Filepath, d.rotateNum)

	return d.rotatingFile(newFilepath)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (d *RotateDispatcher) rotatingFile(newFilepath string) error {
	// close the current file
	if err := d.Close(); err != nil {
		return err
	}

	// rename current to new file.
	err := os.Rename(d.Filepath, newFilepath)
	if err != nil {
		return err
	}

	// reopen file
	d.file, err = fsutil.QuickOpenFile(d.Filepath)
	if err != nil {
		return err
	}

	// reset d.written
	d.written = 0
	return nil
}
