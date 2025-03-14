package rotatefile

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/stdio"
	"github.com/gookit/goutil/strutil"
)

// Writer a flush, close, writer and support rotate file.
//
// refer https://github.com/flike/golog/blob/master/filehandler.go
type Writer struct {
	// writer instance id, use for debug
	id string
	mu sync.RWMutex
	// config of the writer
	cfg *Config

	// current opened logfile
	file *os.File
	// current opened file path. NOTE it maybe not equals Config.Filepath
	path string
	// The original file dir path for the Config.Filepath
	fileDir string
	// The original name and ext information
	fileName, onlyName, fileExt string

	// logfile max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration
	// oldFiles []string
	cleanCh chan struct{}
	stopCh  chan struct{}

	// context use for rotating file by size
	written   uint64 // written size
	rotateNum uint   // rotate times number

	// ---- context use for rotating file by time ----

	// the rotating file name suffix format. eg: "20210102", "20210102_1500"
	suffixFormat   string
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
	// dirSep := filepath.Separator
	// d.fileDir = filepath.Dir(logfile)
	d.fileDir, d.fileName = filepath.Split(d.cfg.Filepath)
	d.fileExt = filepath.Ext(d.fileName)                   // eg: .log
	d.onlyName = strings.TrimSuffix(d.fileName, d.fileExt) // eg: error
	// removes the trailing separator on the dir path
	if ln := len(d.fileDir); ln > 1 && d.fileDir[ln-1] == filepath.Separator {
		d.fileDir = d.fileDir[:ln-1]
	}

	d.backupDur = d.cfg.backupDuration()
	d.suffixFormat = d.cfg.RotateTime.TimeFormat()
	d.checkInterval = d.cfg.RotateTime.Interval()

	// calc and storage next rotating time
	if d.checkInterval > 0 {
		now := d.cfg.TimeClock.Now()
		// next rotating time
		d.nextRotatingAt = d.cfg.RotateTime.FirstCheckTime(now)
		if d.cfg.RotateMode == ModeCreate {
			// logfile = d.cfg.Filepath + "." + now.Format(d.suffixFormat)
			logfile = d.buildFilePath(now.Format(d.suffixFormat))
		}
	}

	// open the current file
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

// Write data to file. then check and do rotate file, then async clean backups
func (d *Writer) Write(p []byte) (n int, err error) {
	// do write data
	if n, err = d.doWrite(p); err != nil {
		return
	}

	// do rotate file
	err = d.doRotate()
	// async clean backup files.
	if err == nil && d.shouldClean(true) {
		d.asyncClean()
	}
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
func (d *Writer) Rotate() error {
	err := d.doRotate()

	// async clean backup files.
	if err == nil && d.shouldClean(true) {
		d.asyncClean()
	}
	return err
}

// do rotate the logfile by config
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
	return
}

// TIP: should only call on d.checkInterval > 0
func (d *Writer) rotatingByTime() error {
	now := d.cfg.TimeClock.Now()
	if now.Before(d.nextRotatingAt) {
		return nil
	}

	// generate new file path.
	// eg: /tmp/error.log => /tmp/error.20220423_1600.log
	// file := d.cfg.Filepath + "." + d.nextRotatingAt.Format(d.suffixFormat)
	file := d.buildFilePath(d.nextRotatingAt.Format(d.suffixFormat))
	err := d.rotatingFile(file, false)

	// calc and storage next rotating time
	d.nextRotatingAt = d.nextRotatingAt.Add(time.Duration(d.checkInterval) * time.Second)
	return err
}

func (d *Writer) rotatingBySize() error {
	d.rotateNum++
	now := d.cfg.TimeClock.Now()
	// up: use now minutes + seconds as rotate number
	numStr := fmt.Sprintf("%d%d%d", now.Hour(), now.Minute(), now.Second())
	numInt := strutil.IntOr(numStr, 0) + now.Nanosecond()/1000
	rotateNum := uint(numInt) + d.rotateNum

	var bakFile string
	if d.cfg.IsMode(ModeCreate) {
		// eg: /tmp/error.log => /tmp/error.894136.log
		// eg: /tmp/error.20220423_1600.log => /tmp/error.20220423_1600_894136.log
		pathNoExt := d.path[:len(d.path)-len(d.fileExt)]
		bakFile = fmt.Sprintf("%s_%d%s", pathNoExt, rotateNum, d.fileExt)
	} else if d.cfg.RenameFunc != nil {
		// rename current to new file by custom RenameFunc
		// eg: /tmp/error.log => /tmp/error.163021_894136.log
		bakFile = d.cfg.RenameFunc(d.cfg.Filepath, rotateNum)
	} else {
		// eg: /tmp/error.log => /tmp/error.25031615_894136.log
		bakFile = d.buildFilePath(fmt.Sprintf("%s_%d", now.Format("06010215"), rotateNum))
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
	d.mu.RLock()
	defer d.mu.RUnlock()

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
				printErrln("rotatefile: clean old files error:", d.doClean())
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

	// up: 单独运行清理，不需要设置 skipSeconds
	return d.doClean(0)
}

// do clean old files by config
//
// - skipSeconds: skip find files that are within the specified seconds
func (d *Writer) doClean(skipSeconds ...int) (err error) {
	// oldFiles: xx.log.yy files, no gz file
	var oldFiles, gzFiles []fileInfo
	fileDir, fileName := d.fileDir, d.fileName
	curFileName := filepath.Base(d.path)

	// FIX: do not process recent changes to avoid conflicts
	skipSec := 30
	if len(skipSeconds) > 0 {
		skipSec = skipSeconds[0]
	}
	limitTime := d.cfg.TimeClock.Now().Add(-time.Second * time.Duration(skipSec))

	// find and clean old files
	d.debugLog("Clean - find old files, match name:", fileName, ", in dir:", fileDir)
	err = fsutil.FindInDir(fileDir, func(fPath string, ent fs.DirEntry) error {
		fi, err := ent.Info()
		if err != nil {
			return err
		}

		// fix: exclude the current file
		if ent.Name() == curFileName {
			return nil
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
	remNum := mathutil.Max(gzNum+oldNum-int(d.cfg.BackupNum), 0)
	d.debugLog("clean old files, gzNum:", gzNum, "oldNum:", oldNum, "remNum:", remNum)

	if remNum > 0 && d.cfg.BackupNum > 0 {
		// remove old gz files
		if gzNum > 0 {
			remNum, err = d.removeOldGzFiles(remNum, gzFiles)
			if err != nil {
				return err
			}
		}

		// remove old log files
		if remNum > 0 && oldNum > 0 {
			oldFiles, err = d.removeOldFiles(remNum, oldFiles)
			if err != nil {
				return err
			}
		}
	}

	if d.cfg.Compress && len(oldFiles) > 0 {
		d.debugLog("compress old normal files to gz files")
		err = d.compressFiles(oldFiles)
	}
	return
}

// remove old gz files
func (d *Writer) removeOldGzFiles(remNum int, gzFiles []fileInfo) (rn int, err error) {
	gzNum := len(gzFiles)
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
		return remNum, errorx.Wrap(err, "remove old gz file error")
	}
	return remNum, nil
}

// remove old log files
func (d *Writer) removeOldFiles(remNum int, oldFiles []fileInfo) (files []fileInfo, err error) {
	// sort by mod-time, oldest at first.
	sort.Sort(modTimeFInfos(oldFiles))
	d.debugLog("remove old normal files ...")

	var idx int
	oldNum := len(oldFiles)

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
		return nil, errorx.Wrap(err, "remove old file error")
	}
	return oldFiles, nil
}

//
// ---------------------------------------------------------------------------
// helper methods
// ---------------------------------------------------------------------------
//

// open the current file. and set the d.file, d.path
func (d *Writer) openFile(logfile string) error {
	file, err := fsutil.OpenFile(logfile, DefaultFileFlags, d.cfg.FilePerm)
	if err != nil {
		return err
	}

	d.path = logfile
	d.file = file
	return nil
}

// return eg. logs/error.20220423_1600.log
func (d *Writer) buildFilePath(suffix string) string {
	fileName := d.onlyName + "." + suffix + d.fileExt
	return fmt.Sprintf("%s/%s", d.fileDir, fileName)
}

func (d *Writer) buildFilterFns(fileName string) []fsutil.FilterFunc {
	onlyName := d.onlyName
	filterFns := []fsutil.FilterFunc{
		fsutil.OnlyFindFile,
		// filter by name. match pattern like: error.log.* eg: error.log.xx, error.log.xx.gz
		func(fPath string, ent fs.DirEntry) bool {
			// ok, _ := path.Match(fileName+".*", ent.Name())
			if !strings.HasPrefix(ent.Name(), fileName) {
				// 自定义文件名 eg: error.log -> error.20220423_02.log
				return strings.HasPrefix(ent.Name(), onlyName)
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
