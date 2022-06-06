package rotatefile

import (
	"os"
	"time"
)

// CConfig struct for clean
type CConfig struct {
	// BackupNum max number for keep old files.
	// 0 is not limit, default is 20.
	BackupNum uint `json:"backup_num" yaml:"backup_num"`

	// BackupTime max time for keep old files.
	// 0 is not limit, default is a week.
	//
	// unit is hours
	BackupTime uint `json:"backup_time" yaml:"backup_time"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`

	// TimeClock for rotate
	TimeClock Clocker
}

// FilesClear multi files by time. TODO
// use for rotate and clear other program produce log files
type FilesClear struct {
	// mu  sync.Mutex
	cfg *CConfig

	namePattern  string
	filepathDirs []string
	// full file path patterns
	filePatterns []string
	// file max backup time. equals Config.BackupTime * time.Hour
	backupDur time.Duration
	// skip error
	skipError bool
}

// NewCleanFiles instance
func NewCleanFiles(cfg *CConfig) *FilesClear {
	return &FilesClear{
		cfg: cfg,
	}
}

//
// ---------------------------------------------------------------------------
// clean backup files
// ---------------------------------------------------------------------------
//

// async clean old files by config
func (r *FilesClear) asyncCleanBackups() {
	if r.cfg.BackupNum == 0 && r.cfg.BackupTime == 0 {
		return
	}

	go func() {
		err := r.Clean()
		if err != nil {
			printErrln("files-clear: clean backup files error:", err)
		}
	}()
}

// Clean old files by config
func (r *FilesClear) Clean() (err error) {
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

func (r *FilesClear) cleanByBackupNum(filePattern string, bckNum int) (err error) {
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

func (r *FilesClear) cleanByBackupTime(filePattern string, cutTime time.Time) (err error) {
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
