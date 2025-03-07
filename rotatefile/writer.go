package rotatefile

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/stdio"
)

// Writer a flush, close, writer and support rotate file.
//
// refer https://github.com/flike/golog/blob/master/filehandler.go
type Writer struct {
	// writer instance id
	id string
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
	suffixFormat   string    // the rotating file name suffix. eg: "20210102", "20210102_1500"
	checkInterval  int64     // check interval seconds.
	nextRotatingAt time.Time // next rotating time
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
	d.id = fmt.Sprintf("%p", d)

	logfile := d.cfg.Filepath
	d.fileDir = filepath.Dir(logfile)
	d.backupDur = d.cfg.backupDuration()

	// if d.cfg.BackupNum > 0 {
	// 	d.oldFiles = make([]string, 0, int(float32(d.cfg.BackupNum)*1.6))
	// }

	d.suffixFormat = d.cfg.RotateTime.TimeFormat()
	d.checkInterval = d.cfg.RotateTime.Interval()

	// calc and storage next rotating time
	if d.checkInterval > 0 {
		now := d.cfg.TimeClock.Now()
		// next rotating time
		d.nextRotatingAt = d.cfg.RotateTime.FirstCheckTime(now)
		if d.cfg.RotateMode == ModeCreate {
			logfile = d.cfg.Filepath + "." + now.Format(d.suffixFormat)
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

// MustClose the writer. alias of Close(), but will panic if has error.
func (d *Writer) MustClose() {
	printErrln("close writer -", d.close(true))
}

func (d *Writer) close(closeStopCh bool) error {
	if err := d.file.Sync(); err != nil {
		return err
	}

	// stop the async clean backups
	if closeStopCh && d.stopCh != nil {
		d.debugLog("close stopCh for stop async clean old files")
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
	// do write data
	if n, err = d.doWrite(p); err != nil {
		return
	}

	// do rotate file
	err = d.doRotate()
	return
}

func (d *Writer) doWrite(p []byte) (n int, err error) {
	// if enable lock
	if !d.cfg.CloseLock {
		d.mu.Lock()
		defer d.mu.Unlock()
	}

	n, err = d.file.Write(p)
	if err == nil {
		// update size
		d.written += uint64(n)
	}
	return
}

// Rotate the file by config and async clean backups
func (d *Writer) Rotate() error { return d.doRotate() }

// do rotate the logfile by config and async clean backups
func (d *Writer) doRotate() (err error) {
	// if enable lock
	if !d.cfg.CloseLock {
		d.mu.Lock()
		defer d.mu.Unlock()
	}

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

	// async clean backup files.
	if d.shouldClean(true) {
		d.asyncClean()
	}
	return
}

// TIP: should only call on d.checkInterval > 0
func (d *Writer) rotatingByTime() error {
	now := d.cfg.TimeClock.Now()
	if now.Before(d.nextRotatingAt) {
		return nil
	}

	// generate new file path.
	// eg: /tmp/error.log => /tmp/error.log.20220423_1600
	file := d.cfg.Filepath + "." + d.nextRotatingAt.Format(d.suffixFormat)
	err := d.rotatingFile(file, false)

	// calc and storage next rotating time
	d.nextRotatingAt = d.nextRotatingAt.Add(time.Duration(d.checkInterval) * time.Second)
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

//
// ---------------------------------------------------------------------------
// clean backup files
// ---------------------------------------------------------------------------
//

// check should clean old files by config
func (d *Writer) shouldClean(withRand bool) bool {
	cfgIsYes := d.cfg.BackupNum > 0 || d.cfg.BackupTime > 0
	if !withRand {
		return cfgIsYes
	}

	// 20% probability trigger clean
	return cfgIsYes && rand.Intn(100) < 20
}

// async clean old files by config. should be in lock.
func (d *Writer) asyncClean() {
	if !d.shouldClean(false) {
		return
	}

	// if already running, send a signal
	if d.cleanCh != nil {
		d.notifyClean()
		return
	}

	// add lock for deny concurrent clean
	d.mu.Lock()
	defer d.mu.Unlock()

	// re-check d.cleanCh is not nil
	if d.cleanCh != nil {
		d.notifyClean()
		return
	}

	// init clean channel
	d.debugLog("INIT clean and stop channels for clean old files")
	d.cleanCh = make(chan struct{})
	d.stopCh = make(chan struct{})

	// start a goroutine to clean backups
	go func() {
		d.debugLog("START a goroutine consumer for clean old files")

		// consume the signal until stop
		for {
			select {
			case <-d.cleanCh:
				d.debugLog("receive signal - clean old files handling")
				printErrln("rotatefile: clean old files error:", d.Clean())
			case <-d.stopCh:
				d.cleanCh = nil
				d.debugLog("STOP consumer for clean old files")
				return // stop clean
			}
		}
	}()
}

func (d *Writer) notifyClean() {
	select {
	case d.cleanCh <- struct{}{}: // notify clean old files
		d.debugLog("sent signal - start clean old files...")
	default: // skip on blocking
		d.debugLog("clean old files signal blocked, SKIP")
	}
}

// Clean old files by config
func (d *Writer) Clean() (err error) {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return errorx.Err("clean: backupNum and backupTime are both 0")
	}

	if !d.mu.TryRLock() {
		d.debugLog("Clean - tryLock=false, SKIP clean old files")
		return nil
	}
	defer d.mu.RUnlock()

	// oldFiles: xx.log.yy files, no gz file
	var oldFiles, gzFiles []fileInfo
	fileDir, fileName := filepath.Split(d.cfg.Filepath)
	if len(fileDir) > 0 {
		// removes the trailing separator
		fileDir = fileDir[:len(fileDir)-1]
	}
	// up: do not process recent changes to avoid conflicts
	limitTime := d.cfg.TimeClock.Now().Add(-time.Second * 30)

	// find and clean old files
	d.debugLog("find old files, match name:", fileName, ", in dir:", fileDir)
	err = fsutil.FindInDir(fileDir, func(fPath string, ent fs.DirEntry) error {
		fi, err := ent.Info()
		if err != nil {
			return err
		}

		if strings.HasSuffix(ent.Name(), compressSuffix) {
			gzFiles = append(gzFiles, newFileInfo(fPath, fi))
		} else if fi.ModTime().Before(limitTime) {
			oldFiles = append(oldFiles, newFileInfo(fPath, fi))
		}
		return nil
	}, d.buildFilterFns(fileName)...)

	gzNum := len(gzFiles)
	oldNum := len(oldFiles)
	remNum := gzNum + oldNum - int(d.cfg.BackupNum)
	d.debugLog("clean old files, gzNum:", gzNum, "oldNum:", oldNum, "remNum:", remNum)

	if remNum > 0 {
		// remove old gz files
		if gzNum > 0 {
			sort.Sort(modTimeFInfos(gzFiles)) // sort by mod-time
			d.debugLog("remove old gz files ...")

			for idx := 0; idx < gzNum; idx++ {
				d.debugLog("remove old gz file:", gzFiles[idx].filePath)
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
			d.debugLog("remove old normal files ...")

			var idx int
			for idx = 0; idx < oldNum; idx++ {
				d.debugLog("remove old file:", oldFiles[idx].filePath)
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
		d.debugLog("compress old normal files to gz files")
		err = d.compressFiles(oldFiles)
	}
	return
}

//
// ---------------------------------------------------------------------------
// helper methods
// ---------------------------------------------------------------------------
//

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

func (d *Writer) buildFilterFns(fileName string) []fsutil.FilterFunc {
	nameNoSuffix := strings.TrimSuffix(fileName, path.Ext(fileName))
	filterFns := []fsutil.FilterFunc{
		fsutil.OnlyFindFile,
		// filter by name. match pattern like: error.log.* eg: error.log.xx, error.log.xx.gz
		func(fPath string, ent fs.DirEntry) bool {
			// ok, _ := path.Match(fileName+".*", ent.Name())
			if !strings.HasPrefix(ent.Name(), fileName) {
				// 自定义文件名 eg: error.log -> error.20220423_02.log
				return strings.HasPrefix(ent.Name(), nameNoSuffix)
			}
			return true
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
			d.debugLog("remove expired file:", fPath)
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
		d.debugLog("compress and rm old file:", fi.filePath)
		if err = os.Remove(fi.filePath); err != nil {
			return errorx.Wrap(err, "remove file error after compress")
		}
	}
	return nil
}

// Debug print debug message on development
func (d *Writer) debugLog(vs ...any) {
	if d.cfg.DebugMode {
		stdio.WriteString("[rotatefile.DEBUG] ID:" + d.id + " | " + fmt.Sprintln(vs...))
	}
}
