package rotatefile_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog/internal"
	"github.com/gookit/slog/rotatefile"
)

func TestNewWriter(t *testing.T) {
	testFile := "testdata/test_writer.log"
	assert.NoErr(t, fsutil.DeleteIfExist(testFile))

	w, err := rotatefile.NewConfig(testFile).Create()
	assert.NoErr(t, err)

	c := w.Config()
	// dump.P(c)
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

func TestWriter_Rotate_modeCreate(t *testing.T) {
	logfile := "testdata/mode_create.log"

	c := rotatefile.NewConfig(logfile)
	c.RotateMode = rotatefile.ModeCreate

	wr, err := c.Create()
	assert.NoErr(t, err)
	_, err = wr.WriteString("[INFO] this is a log message\n")
	assert.NoErr(t, err)
	assert.False(t, fsutil.IsFile(logfile))

	ls, err := filepath.Glob("testdata/mode_create*")
	assert.NoErr(t, err)
	assert.Len(t, ls, 1)

	for i := 0; i < 20; i++ {
		_, err = wr.WriteString("[INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
	}

	// test clean and backup
	c.BackupNum = 2
	c.MaxSize = 128
	err = wr.Rotate()
	assert.NoErr(t, err)
	_, err = wr.WriteString("hi, rotated\n")
	assert.NoErr(t, err)
}

func TestWriter_rotateByTime(t *testing.T) {
	logfile := "testdata/rotate-by-time.log"
	c := rotatefile.NewConfig(logfile).With(func(c *rotatefile.Config) {
		c.DebugMode = true
		c.Compress = true
		c.RotateTime = rotatefile.EverySecond * 2
	})

	w, err := c.Create()
	assert.NoErr(t, err)
	defer func() {
		_ = w.Close()
	}()

	for i := 0; i < 5; i++ {
		_, err = w.WriteString("[INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
		time.Sleep(time.Second)
	}

	files := fsutil.Glob(internal.BuildGlobPattern(logfile))
	dump.P(files)

}

func TestWriter_Clean(t *testing.T) {
	logfile := "testdata/writer_clean.log"

	c := rotatefile.NewConfig(logfile)
	c.MaxSize = 128 // will rotate by size

	wr, err := c.Create()
	assert.NoErr(t, err)
	defer func() {
		_ = wr.Close()
	}()

	for i := 0; i < 20; i++ {
		_, err = wr.WriteString("[INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
	}

	assert.True(t, fsutil.IsFile(logfile))
	_, err = wr.WriteString("hi\n")
	assert.NoErr(t, err)

	files := fsutil.Glob(internal.BuildGlobPattern(logfile))
	dump.P(files)

	// test clean error
	t.Run("clean error", func(t *testing.T) {
		c.BackupNum = 0
		c.BackupTime = 0
		assert.Err(t, wr.Clean())
	})

	// test clean and compress backup
	t.Run("clean and compress", func(t *testing.T) {
		c.BackupNum = 2
		c.Compress = true
		err = wr.Clean()
		assert.NoErr(t, err)

		files := fsutil.Glob(internal.BuildGlobPattern(logfile))
		assert.Lt(t, 2, len(files))
	})
}

// test writer compress
func TestWriter_Compress(t *testing.T) {
	logfile := "testdata/test_compress.log"

	c := rotatefile.NewConfig(logfile)
	c.MaxSize = 128 // will rotate by size
	c.With(rotatefile.WithDebugMode)

	wr, err := c.Create()
	assert.NoErr(t, err)

	for i := 0; i < 20; i++ {
		_, err = wr.WriteString("[INFO] this is a log message, idx=" + mathutil.String(i) + "\n")
		assert.NoErr(t, err)
	}

	assert.True(t, fsutil.IsFile(logfile))
	_, err = wr.WriteString("hi\n")
	assert.NoErr(t, err)
	wr.MustClose()

	files := fsutil.Glob(internal.BuildGlobPattern(logfile))
	assert.NotEmpty(t, files)
	dump.P(files)

	// test clean and compress backup
	t.Run("compress backup", func(t *testing.T) {
		c := rotatefile.NewConfig(logfile,
			rotatefile.WithDebugMode, rotatefile.WithCompress,
			rotatefile.WithBackupNum(2),
		)

		wr, err := c.Create()
		assert.NoErr(t, err)
		defer wr.MustClose()

		err = wr.Clean()
		assert.NoErr(t, err)

		files := fsutil.Glob(internal.BuildGlobPattern(logfile))
		assert.Lt(t, 2, len(files))
		dump.P(files)
	})
}

// TODO set github.com/benbjohnson/clock for mock clock
type constantClock time.Time

func (c constantClock) Now() time.Time { return time.Time(c) }
func (c constantClock) NewTicker(d time.Duration) *time.Ticker {
	return &time.Ticker{}
}
