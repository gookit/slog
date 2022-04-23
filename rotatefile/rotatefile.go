package rotatefile

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/gookit/goutil/fsutil"
)

// RotateFiles multi files. TODO
// use for rotate and clear other program produce logs
type RotateFiles struct {
	patterns string
}

// RotateDispatcher dispatcher for rotate file
//
// refer file-rotatelogs
// refer https://github.com/flike/golog/blob/master/filehandler.go
type RotateDispatcher struct {
	sync.Mutex
	cfg  *Config
	file *os.File
	// file dir path for the Config.Filepath
	fileDir  string
	oldFiles []string
	// file max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration

	// context use for rotating file by size
	written   uint64 // written size
	maxSize   uint64 // file max size byte. equals Config.MaxSize * oneMByte
	rotateNum uint   // rotate times number

	// context use for rotating file by time
	suffixFormat   string // the rotating file name suffix.
	checkInterval  int64
	nextRotatingAt int64
}

// New create rotate dispatcher with config.
func New(c *Config) (*RotateDispatcher, error) {
	d := &RotateDispatcher{
		cfg: c,
	}

	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// NewFromPath create rotate dispatcher with file path.
func NewFromPath(filePath string) (*RotateDispatcher, error) {
	return New(NewConfig(filePath))
}

// init rotate dispatcher
func (d *RotateDispatcher) init() error {
	d.fileDir = path.Dir(d.cfg.Filepath)
	d.maxSize = d.cfg.maxSizeByte()
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

// ReopenFile the log file
func (d *RotateDispatcher) ReopenFile() error {
	if d.file != nil {
		d.file.Close()
	}

	return d.openFile()
}

// ReopenFile the log file
func (d *RotateDispatcher) openFile() error {
	file, err := fsutil.OpenFile(d.cfg.Filepath, DefaultFileFlags, DefaultFilePerm)
	if err != nil {
		return err
	}

	d.file = file
	return nil
}

// Flush sync data to disk. alias of Sync()
func (d *RotateDispatcher) Flush() error {
	return d.file.Sync()
}

// Sync data to disk.
func (d *RotateDispatcher) Sync() error {
	return d.file.Sync()
}

// Close the dispatcher.
// will sync data to disk, then close the file handle
func (d *RotateDispatcher) Close() error {
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
func (d *RotateDispatcher) asyncCleanBackups() {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return
	}

	go func() {
		err := d.CleanBackups()
		if err != nil {
			fmt.Println("slog: async clean backup files error:", err)
		}
	}()
}

// CleanBackups clean old files by config
func (d *RotateDispatcher) CleanBackups() (err error) {
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
func (d *RotateDispatcher) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Write data to file. then check and do rotate file.
func (d *RotateDispatcher) Write(p []byte) (n int, err error) {
	// if enable lock
	if d.cfg.CloseLock == false {
		d.Lock()
		defer d.Unlock()
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
func (d *RotateDispatcher) Rotate() (err error) {
	// do rotate file by time
	if d.checkInterval > 0 {
		err = d.rotatingByTime()
		if err != nil {
			return
		}
	}

	// do rotate file by size
	if d.cfg.MaxSize > 0 && d.written >= d.maxSize {
		err = d.rotatingBySize()
	}
	return
}

func (d *RotateDispatcher) rotatingByTime() error {
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

func (d *RotateDispatcher) rotatingBySize() error {
	// rename current to new file
	d.rotateNum++

	// eg: /tmp/error.log => /tmp/error.log.163021_1
	newFilepath := d.cfg.RenameFunc(d.cfg.Filepath, d.rotateNum)

	return d.rotatingFile(newFilepath)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (d *RotateDispatcher) rotatingFile(newFilepath string) error {
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
