package rotatefile_test

import (
	"log"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog/rotatefile"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := rotatefile.NewConfig("testdata/test.log")

	assert.Equal(t, rotatefile.DefaultBackNum, cfg.BackupNum)
	assert.Equal(t, rotatefile.DefaultBackTime, cfg.BackupTime)
	assert.Equal(t, rotatefile.EveryHour, cfg.RotateTime)
	assert.Equal(t, rotatefile.DefaultMaxSize, cfg.MaxSize)

	dump.P(cfg)
}

func TestNewWriter(t *testing.T) {
	testFile := "testdata/test.log"
	assert.NoError(t, fsutil.DeleteIfExist(testFile))

	wr, err := rotatefile.NewConfig(testFile).Create()
	if err != nil {
		return
	}

	c := wr.Config()
	assert.Equal(t, c.MaxSize, rotatefile.DefaultMaxSize)

	w, err := c.Create()
	assert.NoError(t, err)

	_, err = w.WriteString("info log message\n")
	assert.NoError(t, err)
	assert.True(t, fsutil.IsFile(testFile))

	assert.NoError(t, w.Flush())
	assert.NoError(t, w.Close())
}

func ExampleNewWriter_on_other_logger() {
	logFile := "testdata/another_logger.log"
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err)
	}

	log.SetOutput(writer)
	log.Println("log message")
}
