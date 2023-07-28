package rotatefile

import (
	"compress/gzip"
	"io"
	"io/fs"
	"os"

	"github.com/gookit/goutil/fsutil"
)

const compressSuffix = ".gz"

func compressFile(srcPath, dstPath string) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create and open a gz file
	gzFile, err := fsutil.OpenTruncFile(dstPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	srcSt, err := srcFile.Stat()
	if err != nil {
		return err
	}

	zw := gzip.NewWriter(gzFile)
	zw.Name = srcSt.Name()
	zw.ModTime = srcSt.ModTime()

	// do copy
	if _, err = io.Copy(zw, srcFile); err != nil {
		_ = zw.Close()
		return err
	}
	return zw.Close()
}

// TODO replace to fsutil.FileInfo
type fileInfo struct {
	fs.FileInfo
	filePath string
}

// Path get file full path. eg: "/path/to/file.go"
func (fi *fileInfo) Path() string {
	return fi.filePath
}

func newFileInfo(filePath string, fi fs.FileInfo) fileInfo {
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
