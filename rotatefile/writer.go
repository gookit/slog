package rotatefile

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
)

// Writer a flush, close, writer and support rotate file.
//
// refer https://github.com/flike/golog/blob/master/filehandler.go
type Writer struct {
	mu sync.Mutex
	// config of the writer
	cfg *Config
	// current opened logfile
	file *os.File
	path string
	// logfile dir path for the Config.Filepath
	fileDir string

	// logfile max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration
	// oldFiles []string
	cleanCh chan struct{}
	stopCh  chan struct{}

	// context use for rotating file by size
	written   uint64 // written size
	rotateNum uint   // rotate times number

	// context use for rotating file by time
	suffixFormat   string // the rotating file name suffix. eg: "20210102", "20210102_1500"
	checkInterval  int64  // check interval seconds.
	nextRotatingAt int64
}

// NewWriter create rotate write with config and init it.
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
	logfile := d.cfg.Filepath
	d.fileDir = path.Dir(logfile)
	d.backupDur = d.cfg.backupDuration()

	// if d.cfg.BackupNum > 0 {
	// 	d.oldFiles = make([]string, 0, int(float32(d.cfg.BackupNum)*1.6))
	// }

	d.suffixFormat = d.cfg.RotateTime.TimeFormat()
	d.checkInterval = d.cfg.RotateTime.Interval()

	// calc and storage next rotating time
	if d.checkInterval > 0 {
		nowTime := d.cfg.TimeClock.Now()
		// next rotating time
		d.nextRotatingAt = d.cfg.RotateTime.FirstCheckTime(nowTime)

		if d.cfg.RotateMode == ModeCreate {
			logfile = d.cfg.Filepath + "." + nowTime.Format(d.suffixFormat)
		}
	}

	// open the logfile
	return d.openFile(logfile)
}

// Config get the config
func (d *Writer) Config() Config {
	return *d.cfg
}

// Flush sync data to disk. alias of Sync()
func (d *Writer) Flush() error {
	return d.file.Sync()
}

// Sync data to disk.
func (d *Writer) Sync() error {
	return d.file.Sync()
}

// Close the writer. will sync data to disk, then close the file handle.
// and will stop the async clean backups.
func (d *Writer) Close() error {
	return d.close(true)
}

func (d *Writer) close(closeStopCh bool) error {
	if err := d.file.Sync(); err != nil {
		return err
	}

	// stop the async clean backups
	if closeStopCh && d.stopCh != nil {
		d.cfg.Debug("close stopCh for stop async clean old files")
		close(d.stopCh)
		d.stopCh = nil
	}
	return d.file.Close()
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
	if !d.cfg.CloseLock {
		d.mu.Lock()
		defer d.mu.Unlock()
	}

	n, err = d.file.Write(p)
	if err != nil {
		return
	}

	// update written size
	d.written += uint64(n)

	// rotate file
	err = d.doRotate()
	return
}

// Rotate the file by config and async clean backups
func (d *Writer) Rotate() error { return d.doRotate() }

// do rotate the logfile by config and async clean backups
func (d *Writer) doRotate() (err error) {
	// do rotate file by size
	if d.cfg.MaxSize > 0 && d.written >= d.cfg.MaxSize {
		err = d.rotatingBySize()
		if err != nil {
			return
		}
	}

	// do rotate file by time
	if d.checkInterval > 0 && d.written > 0 {
		err = d.rotatingByTime()
	}

	// async clean backup files. TODO only call on file rotated.
	d.asyncClean()
	return
}

// TIP: should only call on d.checkInterval > 0
func (d *Writer) rotatingByTime() error {
	now := d.cfg.TimeClock.Now()
	if d.nextRotatingAt > now.Unix() {
		return nil
	}

	// generate new file path.
	// eg: /tmp/error.log => /tmp/error.log.20220423_1600
	file := d.cfg.Filepath + "." + now.Format(d.suffixFormat)
	err := d.rotatingFile(file, false)

	// storage next rotating time
	d.nextRotatingAt = now.Unix() + d.checkInterval
	return err
}

func (d *Writer) rotatingBySize() error {
	d.rotateNum++

	var bakFile string
	if d.cfg.IsMode(ModeCreate) {
		// eg: /tmp/error.log.20220423_1600 => /tmp/error.log.20220423_1600_001
		bakFile = fmt.Sprintf("%s_%03d", d.path, d.rotateNum)
	} else {
		// rename current to new file
		// eg: /tmp/error.log => /tmp/error.log.163021_001
		bakFile = d.cfg.RenameFunc(d.cfg.Filepath, d.rotateNum)
	}

	// always rename current to new file
	return d.rotatingFile(bakFile, true)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (d *Writer) rotatingFile(bakFile string, rename bool) error {
	// close the current file
	if err := d.close(false); err != nil {
		return err
	}

	// record old files for clean.
	// d.oldFiles = append(d.oldFiles, bakFile)

	// rename current to new file.
	if rename || d.cfg.RotateMode == ModeRename {
		if err := os.Rename(d.path, bakFile); err != nil {
			return err
		}
	}

	// filepath for reopen
	logfile := d.path
	if d.cfg.RotateMode == ModeRename {
		logfile = d.cfg.Filepath
	}

	// reopen log file
	if err := d.openFile(logfile); err != nil {
		return err
	}

	// reset written
	d.written = 0
	return nil
}

// open the log file. and set the d.file, d.path
func (d *Writer) openFile(logfile string) error {
	file, err := fsutil.OpenFile(logfile, DefaultFileFlags, d.cfg.FilePerm)
	if err != nil {
		return err
	}

	d.path = logfile
	d.file = file
	return nil
}

//
// ---------------------------------------------------------------------------
// clean backup files
// ---------------------------------------------------------------------------
//

// async clean old files by config. should be in lock.
func (d *Writer) asyncClean() {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return
	}

	// if already running, send a signal
	if d.cleanCh != nil {
		select {
		case d.cleanCh <- struct{}{}:
			d.cfg.Debug("signal sent start clean old files ")
		default: // skip on blocking
			d.cfg.Debug("clean old files signal blocked, skip")
		}
		return
	}

	// init clean channel
	d.cfg.Debug("init clean/stop channel for clean old files")
	d.cleanCh = make(chan struct{})
	d.stopCh = make(chan struct{})

	// start a goroutine to clean backups
	go func() {
		d.cfg.Debug("start a goroutine consumer for clean old files")

		// consume the signal until stop
		for {
			select {
			case <-d.cleanCh:
				d.cfg.Debug("clean old files handling ...")
				printErrln("rotatefile: clean old files error:", d.Clean())
			case <-d.stopCh:
				d.cleanCh = nil
				d.cfg.Debug("stop consumer for clean old files")
				return // stop clean
			}
		}
	}()
}

// Clean old files by config
func (d *Writer) Clean() (err error) {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return errorx.Err("clean: backupNum and backupTime are both 0")
	}

	// oldFiles: xx.log.yy files, no gz file
	var oldFiles, gzFiles []fileInfo
	fileDir, fileName := path.Split(d.cfg.Filepath)

	// find and clean old files
	err = fsutil.FindInDir(fileDir, func(fPath string, ent fs.DirEntry) error {
		fi, err := ent.Info()
		if err != nil {
			return err
		}

		if strings.HasSuffix(ent.Name(), compressSuffix) {
			gzFiles = append(gzFiles, newFileInfo(fPath, fi))
		} else {
			oldFiles = append(oldFiles, newFileInfo(fPath, fi))
		}
		return nil
	}, d.buildFilterFns(fileName)...)

	gzNum := len(gzFiles)
	oldNum := len(oldFiles)
	remNum := gzNum + oldNum - int(d.cfg.BackupNum)
	d.cfg.Debug("clean old files, gzNum:", gzNum, "oldNum:", oldNum, "remNum:", remNum)

	if remNum > 0 {
		// remove old gz files
		if gzNum > 0 {
			sort.Sort(modTimeFInfos(gzFiles)) // sort by mod-time
			d.cfg.Debug("remove old gz files ...")

			for idx := 0; idx < gzNum; idx++ {
				if err = os.Remove(gzFiles[idx].filePath); err != nil {
					break
				}

				remNum--
				if remNum == 0 {
					break
				}
			}

			if err != nil {
				return errorx.Wrap(err, "remove old gz file error")
			}
		}

		// remove old log files
		if remNum > 0 && oldNum > 0 {
			// sort by mod-time, oldest at first.
			sort.Sort(modTimeFInfos(oldFiles))
			d.cfg.Debug("remove old normal files ...")

			var idx int
			for idx = 0; idx < oldNum; idx++ {
				if err = os.Remove(oldFiles[idx].filePath); err != nil {
					break
				}

				remNum--
				if remNum == 0 {
					break
				}
			}

			oldFiles = oldFiles[idx+1:]
			if err != nil {
				return errorx.Wrap(err, "remove old file error")
			}
		}
	}

	if d.cfg.Compress && len(oldFiles) > 0 {
		d.cfg.Debug("compress old normal files to gz files")
		err = d.compressFiles(oldFiles)
	}
	return
}

func (d *Writer) buildFilterFns(fileName string) []fsutil.FilterFunc {
	filterFns := []fsutil.FilterFunc{
		fsutil.OnlyFindFile,
		// filter by name. match pattern like: error.log.*
		// eg: error.log.xx, error.log.xx.gz
		func(fPath string, ent fs.DirEntry) bool {
			ok, _ := path.Match(fileName+".*", ent.Name())
			return ok
		},
	}

	// filter by mod-time, clear expired files
	if d.cfg.BackupTime > 0 {
		cutTime := d.cfg.TimeClock.Now().Add(-d.backupDur)
		filterFns = append(filterFns, func(fPath string, ent fs.DirEntry) bool {
			fi, err := ent.Info()
			if err != nil {
				return false // skip, not handle
			}

			// collect un-expired
			if fi.ModTime().After(cutTime) {
				return true
			}

			// remove expired files
			printErrln("rotatefile: remove expired file error:", os.Remove(fPath))
			return false
		})
	}

	return filterFns
}

func (d *Writer) compressFiles(oldFiles []fileInfo) error {
	for _, fi := range oldFiles {
		err := compressFile(fi.filePath, fi.filePath+compressSuffix)
		if err != nil {
			return errorx.Wrap(err, "compress old file error")
		}

		// remove old log file
		if err = os.Remove(fi.filePath); err != nil {
			return errorx.Wrap(err, "remove file error after compress")
		}
	}
	return nil
}
