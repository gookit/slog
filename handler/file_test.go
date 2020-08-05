package handler_test

import (
	"io/ioutil"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

const testFile = "./testdata/app.log"
const testSubFile = "./testdata/subdir/app.log"

func TestNewFileHandler(t *testing.T) {
	assert.NoError(t, fsutil.DeleteIfFileExist(testFile))
	h := handler.NewFileHandler(testFile, false).Configure(func(h *handler.FileHandler) {
		h.NoBuffer = true
	})

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
