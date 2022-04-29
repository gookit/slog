package rotatefile

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RotateFiles multi files by time. TODO
// use for rotate and clear other program produce log files
//
// refer file-rotatelogs
type RotateFiles struct {
	mu  sync.Mutex
	cfg *Config

	namePattern  string
	filepathDirs []string
	// full file path patterns
	filePatterns []string
	// file max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration
	//
	skipError bool
}

// Rotate do rotate handle
func (r *RotateFiles) Rotate() error {
	return nil
}

//
// ---------------------------------------------------------------------------
// clean backup files
// ---------------------------------------------------------------------------
//

// async clean old files by config
func (r *RotateFiles) asyncCleanBackups() {
	if r.cfg.BackupNum == 0 && r.cfg.BackupTime == 0 {
		return
	}

	go func() {
		err := r.Clean()
		if err != nil {
			printlnStderr("rotatefile: clean backup files error:", err)
		}
	}()
}

// Clean old files by config
func (r *RotateFiles) Clean() (err error) {
	// clear by time, can also clean by number
	if r.cfg.BackupTime > 0 {
		cutTime := r.cfg.TimeClock.Now().Add(-r.backupDur)
		for _, fileDir := range r.filepathDirs {
			// eg: /tmp/ + error.log.* => /tmp/error.log.*
			filePattern := fileDir + r.namePattern

			err = r.cleanByBackupTime(filePattern, cutTime)
			if err != nil {
				return err
			}
		}

		for _, pattern := range r.filePatterns {
			err = r.cleanByBackupTime(pattern, cutTime)
			if err != nil {
				break
			}
		}
		return
	}

	if r.cfg.BackupNum == 0 {
		return nil
	}

	// clear by number.
	bckNum := int(r.cfg.BackupNum)
	for _, fileDir := range r.filepathDirs {
		pattern := fileDir + r.namePattern

		err = r.cleanByBackupNum(pattern, bckNum)
		if err != nil {
			return err
		}
	}

	for _, pattern := range r.filePatterns {
		err = r.cleanByBackupNum(pattern, bckNum)
		if err != nil {
			break
		}
	}
	return
}

func (r *RotateFiles) cleanByBackupNum(filePattern string, bckNum int) (err error) {
	keepNum := 0
	err = globWithFunc(filePattern, func(oldFile string) error {
		stat, err := os.Stat(oldFile)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return nil
		}

		if keepNum < bckNum {
			keepNum++
		}

		// remove expired files
		return os.Remove(oldFile)
	})

	return
}

func (r *RotateFiles) cleanByBackupTime(filePattern string, cutTime time.Time) (err error) {
	oldFiles := make([]string, 0, 8)

	// match all old rotate files. eg: /tmp/error.log.*
	err = globWithFunc(filePattern, func(filePath string) error {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return nil
		}

		if stat.ModTime().After(cutTime) {
			oldFiles = append(oldFiles, filePath)
			return nil
		}

		// remove expired files
		return os.Remove(filePath)
	})

	// clear by number.
	maxNum := int(r.cfg.BackupNum)
	if maxNum > 0 && len(oldFiles) > maxNum {
		for idx := 0; len(oldFiles) > maxNum; idx++ {
			err = os.Remove(oldFiles[idx])
			if err != nil {
				break
			}
		}
	}

	return
}

func globWithFunc(pattern string, fn func(filePath string) error) (err error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		err = fn(filePath)
		if err != nil {
			break
		}
	}
	return
}
