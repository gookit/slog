package rotatefile

import (
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// Writer a flush, close, writer and support rotate file.
//
// refer https://github.com/flike/golog/blob/master/filehandler.go
type Writer struct {
	mu sync.Mutex
	// config of the writer
	cfg  *Config
	file *os.File
	// file dir path for the Config.Filepath
	fileDir  string
	oldFiles []string
	// file max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration

	// context use for rotating file by size
	written   uint64 // written size
	rotateNum uint   // rotate times number

	// context use for rotating file by time
	suffixFormat   string // the rotating file name suffix.
	checkInterval  int64
	nextRotatingAt int64
}

// NewWriter create rotate dispatcher with config.
func NewWriter(c *Config) (*Writer, error) {
	d := &Writer{cfg: c}

	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// NewWriterWith create rotate writer with some settings.
func NewWriterWith(fns ...ConfigFn) (*Writer, error) {
	return NewWriter(NewConfigWith(fns...))
}

// init rotate dispatcher
func (d *Writer) init() error {
	d.fileDir = path.Dir(d.cfg.Filepath)
	d.backupDur = d.cfg.backupDuration()

	if d.cfg.BackupNum > 0 {
		d.oldFiles = make([]string, 0, int(float32(d.cfg.BackupNum)*1.6))
	}

	// open the log file
	err := d.openFile()
	if err != nil {
		return err
	}

	d.suffixFormat = d.cfg.RotateTime.TimeFormat()
	d.checkInterval = d.cfg.RotateTime.Interval()

	// calc and storage next rotating time
	if d.checkInterval > 0 {
		nowTime := d.cfg.TimeClock.Now()
		d.nextRotatingAt = d.cfg.RotateTime.FirstCheckTime(nowTime)
	}

	return nil
}

// Config get the config
func (d *Writer) Config() *Config {
	return d.cfg
}

// ReopenFile the log file
func (d *Writer) ReopenFile() error {
	if d.file != nil {
		d.file.Close()
	}

	return d.openFile()
}

// ReopenFile the log file
func (d *Writer) openFile() error {
	file, err := fsutil.OpenFile(d.cfg.Filepath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return err
	}

	d.file = file
	return nil
}

// Flush sync data to disk. alias of Sync()
func (d *Writer) Flush() error {
	return d.file.Sync()
}

// Sync data to disk.
func (d *Writer) Sync() error {
	return d.file.Sync()
}

// Close the writer.
// will sync data to disk, then close the file handle
func (d *Writer) Close() error {
	err := d.file.Sync()
	if err != nil {
		return err
	}

	return d.file.Close()
}

//
// ---------------------------------------------------------------------------
// clean backup files
// ---------------------------------------------------------------------------
//

// async clean old files by config
func (d *Writer) asyncCleanBackups() {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return
	}

	go func() {
		err := d.Clean()
		if err != nil {
			printlnStderr("rotatefile: clean backup files error:", err)
		}
	}()
}

// Clean old files by config
func (d *Writer) Clean() (err error) {
	maxNum := int(d.cfg.BackupNum)
	if maxNum > 0 && len(d.oldFiles) > maxNum {
		var idx int
		for idx = 0; len(d.oldFiles) > maxNum; idx++ {
			err = os.Remove(d.oldFiles[idx])
			if err != nil {
				break
			}
		}

		d.oldFiles = d.oldFiles[idx+1:]
		if err != nil {
			return err
		}
	}

	// clear by time
	if d.cfg.BackupTime > 0 {
		// match all old rotate files. eg: /tmp/error.log.*
		files, err := filepath.Glob(d.fileDir + ".*")
		if err != nil {
			return err
		}

		cutTime := d.cfg.TimeClock.Now().Add(-d.backupDur)
		for _, oldFile := range files {
			stat, err := os.Stat(oldFile)
			if err != nil {
				return err
			}

			if stat.ModTime().After(cutTime) {
				continue
			}

			// remove expired files
			err = os.Remove(oldFile)
			if err != nil {
				break
			}
		}
	}

	return
}

//
// ---------------------------------------------------------------------------
// write and rotate file
// ---------------------------------------------------------------------------
//

// WriteString implements the io.StringWriter
func (d *Writer) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Write data to file. then check and do rotate file.
func (d *Writer) Write(p []byte) (n int, err error) {
	// if enable lock
	if d.cfg.CloseLock == false {
		d.mu.Lock()
		defer d.mu.Unlock()
	}

	n, err = d.file.Write(p)
	if err != nil {
		return
	}

	d.written += uint64(n)

	// rotate file
	err = d.Rotate()

	// clean backup files
	d.asyncCleanBackups()
	return
}

// Rotate the file by config.
func (d *Writer) Rotate() (err error) {
	// do rotate file by time
	if d.checkInterval > 0 {
		err = d.rotatingByTime()
		if err != nil {
			return
		}
	}

	// do rotate file by size
	if d.cfg.MaxSize > 0 && d.written >= d.cfg.MaxSize {
		err = d.rotatingBySize()
	}
	return
}

func (d *Writer) rotatingByTime() error {
	now := d.cfg.TimeClock.Now()
	if d.nextRotatingAt > now.Unix() {
		return nil
	}

	// rename current to new file.
	// eg: /tmp/error.log => /tmp/error.log.20220423_1600
	newFilepath := d.cfg.Filepath + "." + now.Format(d.suffixFormat)

	err := d.rotatingFile(newFilepath)

	// storage next rotating time
	d.nextRotatingAt = now.Unix() + d.checkInterval
	return err
}

func (d *Writer) rotatingBySize() error {
	// rename current to new file
	d.rotateNum++

	// eg: /tmp/error.log => /tmp/error.log.163021_1
	newFilepath := d.cfg.RenameFunc(d.cfg.Filepath, d.rotateNum)

	return d.rotatingFile(newFilepath)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (d *Writer) rotatingFile(newFilepath string) error {
	// close the current file
	if err := d.Close(); err != nil {
		return err
	}

	// rename current to new file.
	err := os.Rename(d.cfg.Filepath, newFilepath)
	if err != nil {
		return err
	}

	// record old files for clean.
	d.oldFiles = append(d.oldFiles, newFilepath)

	// reopen file
	if err = d.openFile(); err != nil {
		return err
	}

	// reset written
	d.written = 0
	return nil
}
