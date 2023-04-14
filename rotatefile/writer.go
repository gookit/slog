package rotatefile

import (
	"os"
	"path"
	"path/filepath"
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
	cfg  *Config
	file *os.File
	// file dir path for the Config.Filepath
	fileDir string
	// oldFiles []string
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

	// if d.cfg.BackupNum > 0 {
	// 	d.oldFiles = make([]string, 0, int(float32(d.cfg.BackupNum)*1.6))
	// }

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

	d.written += uint64(n)

	// rotate file
	err = d.doRotate()
	return
}

// Rotate the file by config.
func (d *Writer) Rotate() (err error) {
	return d.doRotate()
}

// Rotate the file by config.
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

	// clean backup files
	d.asyncCleanBackups()
	return
}

func (d *Writer) rotatingByTime() error {
	now := d.cfg.TimeClock.Now()
	if d.nextRotatingAt > now.Unix() {
		return nil
	}

	// rename current to new file.
	// eg: /tmp/error.log => /tmp/error.log.20220423_1600
	bakFilepath := d.cfg.Filepath + "." + now.Format(d.suffixFormat)

	err := d.rotatingFile(bakFilepath)

	// storage next rotating time
	d.nextRotatingAt = now.Unix() + d.checkInterval
	return err
}

func (d *Writer) rotatingBySize() error {
	// rename current to new file
	d.rotateNum++

	// eg: /tmp/error.log => /tmp/error.log.163021_1
	bakFilepath := d.cfg.RenameFunc(d.cfg.Filepath, d.rotateNum)

	return d.rotatingFile(bakFilepath)
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (d *Writer) rotatingFile(bakFilepath string) error {
	// close the current file
	if err := d.Close(); err != nil {
		return err
	}

	// rename current to new file.
	err := os.Rename(d.cfg.Filepath, bakFilepath)
	if err != nil {
		return err
	}

	// record old files for clean.
	// d.oldFiles = append(d.oldFiles, bakFilepath)

	// reopen log file
	if err = d.openFile(); err != nil {
		return err
	}

	// reset written
	d.written = 0
	return nil
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

	// TODO pref: only start once
	go func() {
		err := d.Clean()
		if err != nil {
			printErrln("rotatefile: clean backup files error:", err)
		}
	}()
}

// Clean old files by config
func (d *Writer) Clean() (err error) {
	if d.cfg.BackupNum == 0 && d.cfg.BackupTime == 0 {
		return
	}

	// oldFiles: old xx.log.xx files, no gz file
	var oldFiles, gzFiles []fileInfo
	fileDir, fileName := path.Split(d.cfg.Filepath)

	// find and clean old files
	err = findFilesInDir(fileDir, func(fPath string, fi os.FileInfo) error {
		if strings.HasSuffix(fi.Name(), compressSuffix) {
			gzFiles = append(gzFiles, newFileInfo(fPath, fi))
		} else {
			oldFiles = append(oldFiles, newFileInfo(fPath, fi))
		}

		return nil
	}, d.buildFilterFns(fileName)...)

	gzNum := len(gzFiles)
	oldNum := len(oldFiles)
	maxNum := int(d.cfg.BackupNum)
	rmNum := gzNum + oldNum - maxNum

	if rmNum > 0 {
		// remove old gz files
		if gzNum > 0 {
			sort.Sort(modTimeFInfos(gzFiles)) // sort by mod-time

			for idx := 0; idx < gzNum; idx++ {
				if err = os.Remove(gzFiles[idx].filePath); err != nil {
					break
				}

				rmNum--
				if rmNum == 0 {
					break
				}
			}

			if err != nil {
				return errorx.Wrap(err, "")
			}
		}

		// remove old log files
		if rmNum > 0 && oldNum > 0 {
			sort.Sort(modTimeFInfos(oldFiles)) // sort by mod-time

			var idx int
			for idx = 0; idx < oldNum; idx++ {
				if err = os.Remove(oldFiles[idx].filePath); err != nil {
					break
				}

				rmNum--
				if rmNum == 0 {
					break
				}
			}

			oldFiles = oldFiles[idx+1:]
			if err != nil {
				return err
			}
		}
	}

	if d.cfg.Compress && len(oldFiles) > 0 {
		err = d.compressFiles(oldFiles)
	}
	return
}

func (d *Writer) buildFilterFns(fileName string) []filterFunc {
	filterFns := []filterFunc{
		// filter by name. should match like error.log.*
		// eg: error.log.xx, error.log.xx.gz
		func(fPath string, fi os.FileInfo) bool {
			ok, err := filepath.Match(fileName+".*", fi.Name())
			if err != nil {
				printErrln("rotatefile: match old file error:", err)
				return false // skip, not handle
			}

			return ok
		},
	}

	// filter by mod-time, clear expired files
	if d.cfg.BackupTime > 0 {
		cutTime := d.cfg.TimeClock.Now().Add(-d.backupDur)
		filterFns = append(filterFns, func(fPath string, fi os.FileInfo) bool {
			// collect un-expired
			if fi.ModTime().After(cutTime) {
				return true
			}

			// remove expired files
			err := os.Remove(fPath)
			if err != nil {
				printErrln("rotatefile: remove expired file error:", err)
			}

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

		// remove raw log file
		if err = os.Remove(fi.filePath); err != nil {
			return err
		}
	}
	return nil
}

type fileInfo struct {
	os.FileInfo
	filePath string
}

func newFileInfo(filePath string, fi os.FileInfo) fileInfo {
	return fileInfo{filePath: filePath, FileInfo: fi}
}

// modTimeFInfos sorts by oldest time modified in the fileInfo.
// eg: [old_220211, old_220212, old_220213]
type modTimeFInfos []fileInfo

// Less check
func (fis modTimeFInfos) Less(i, j int) bool {
	return fis[j].ModTime().After(fis[i].ModTime())
}

// Swap value
func (fis modTimeFInfos) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

// Len get
func (fis modTimeFInfos) Len() int {
	return len(fis)
}
