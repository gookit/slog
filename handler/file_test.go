package handler_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

const testFile = "./testdata/app.log"

func deleteIfExist(fpath string) {
	if !fsutil.IsFile(fpath) {
		return
	}

	err := os.Remove(fpath)
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewFileHandler(t *testing.T) {
	deleteIfExist(testFile)
	h := handler.NewFileHandler(testFile, false)

	l := slog.NewWithHandlers(h)
	l.Info("info message")
	l.Warn("warn message")

	assert.True(t, fsutil.IsFile(testFile))
	bts, err := ioutil.ReadFile(testFile)
	assert.NoError(t, err)
	str := string(bts)
	assert.Contains(t, str, "[INFO]")
	assert.Contains(t, str, "info message")
	assert.Contains(t, str, "[WARNING]")
	assert.Contains(t, str, "warn message")
}
