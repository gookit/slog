package rotatefile_test

import (
	"log"
	"testing"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/rotatefile"
)

func ExampleNewWriter_on_other_logger() {
	logFile := "testdata/another_logger.log"
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err)
	}

	log.SetOutput(writer)
	log.Println("log message")
}

func TestNewWriter(t *testing.T) {
	testFile := "testdata/test.log"
	assert.NoErr(t, fsutil.DeleteIfExist(testFile))

	w, err := rotatefile.NewConfig(testFile).Create()
	assert.NoErr(t, err)

	c := w.Config()
	assert.Eq(t, c.MaxSize, rotatefile.DefaultMaxSize)

	_, err = w.WriteString("info log message\n")
	assert.NoErr(t, err)
	assert.True(t, fsutil.IsFile(testFile))

	assert.NoErr(t, w.Sync())
	assert.NoErr(t, w.Flush())
	assert.NoErr(t, w.Close())

	w, err = rotatefile.NewWriterWith(rotatefile.WithFilepath(testFile))
	assert.NoErr(t, err)
	assert.Eq(t, w.Config().Filepath, testFile)
}

func TestWriter_Clean(t *testing.T) {
	logfile := "testdata/test_clean.log"
	assert.NoErr(t, fsutil.DeleteIfExist(logfile))

	c := rotatefile.NewConfig(logfile)
	c.MaxSize = 128
	c.BackupNum = 0
	c.BackupTime = 0

	wr, err := c.Create()
	assert.NoErr(t, err)

	for i := 0; i < 20; i++ {
		_, err = wr.WriteString("[INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
	}

	assert.True(t, fsutil.IsFile(logfile))

	_, err = wr.WriteString("hi\n")
	assert.NoErr(t, err)

	c.BackupNum = 2
	c.Compress = true
	err = wr.Clean()
	assert.NoErr(t, err)
}
