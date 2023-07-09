package rotatefile_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog/rotatefile"
)

func TestFilesClear_Clean(t *testing.T) {
	// make files for clean
	makeNum := 5
	makeWaitCleanFiles("file_clean.log", makeNum)
	_, err := fsutil.PutContents("testdata/subdir/some.txt", "test data")
	assert.NoErr(t, err)

	// create clear
	fc := rotatefile.NewFilesClear()
	fc.WithConfig(rotatefile.NewCConfig())
	fc.WithConfigFn(func(c *rotatefile.CConfig) {
		c.AddDirPath("testdata", "not-exist-dir")
		c.BackupNum = 1
		c.BackupTime = 3
		c.TimeUnit = time.Second // for test
	})

	cfg := fc.Config()
	assert.Eq(t, uint(1), cfg.BackupNum)
	dump.P(cfg)

	// do clean
	assert.NoErr(t, fc.Clean())

	files := fsutil.Glob("testdata/file_clean.log.*")
	dump.P(files)
	assert.NotEmpty(t, files)
	assert.Lt(t, len(files), makeNum)

	t.Run("error", func(t *testing.T) {
		fc := rotatefile.NewFilesClear(func(c *rotatefile.CConfig) {
			c.BackupNum = 0
			c.BackupTime = 0
		})
		assert.Err(t, fc.Clean())
	})
}

func TestFilesClear_DaemonClean(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		fc := rotatefile.NewFilesClear(func(c *rotatefile.CConfig) {
			c.BackupNum = 0
			c.BackupTime = 0
		})
		assert.Panics(t, func() {
			fc.StopDaemon()
		})
		assert.Panics(t, func() {
			fc.DaemonClean(nil)
		})
	})

	fc := rotatefile.NewFilesClear(func(c *rotatefile.CConfig) {
		c.AddPattern("testdata/file_daemon_clean.*")
		c.BackupNum = 1
		c.BackupTime = 3
		c.TimeUnit = time.Second      // for test
		c.CheckInterval = time.Second // for test
	})

	cfg := fc.Config()
	dump.P(cfg)

	// make files for clean
	makeNum := 5
	makeWaitCleanFiles("file_daemon_clean.log", makeNum)

	// test daemon clean
	wg := sync.WaitGroup{}
	wg.Add(1)

	// start daemon
	go fc.DaemonClean(func() {
		fmt.Println("daemon clean stopped, at", timex.Now().DateFormat("ymdTH:i:s.v"))
		wg.Done()
	})

	// stop daemon
	go func() {
		time.Sleep(time.Millisecond * 1200)
		fmt.Println("stop daemon clean, at", timex.Now().DateFormat("ymdTH:i:s.v"))
		fc.StopDaemon()
	}()

	// wait for stop
	wg.Wait()

	files := fsutil.Glob("testdata/file_daemon_clean.log.*")
	dump.P(files)
	assert.NotEmpty(t, files)
	assert.Lt(t, len(files), makeNum)
}

func makeWaitCleanFiles(nameTpl string, makeNum int) {
	for i := 0; i < makeNum; i++ {
		fpath := fmt.Sprintf("testdata/%s.%03d", nameTpl, i)
		fmt.Println("make file:", fpath)
		_, err := fsutil.PutContents(fpath, []byte("test contents ..."))
		goutil.PanicErr(err)
		time.Sleep(time.Second)
	}

	fmt.Println("wait clean files:")
	err := fsutil.GlobWithFunc("./testdata/"+nameTpl+".*", func(fpath string) error {
		fi, err := os.Stat(fpath)
		goutil.PanicErr(err)

		fmt.Printf("  %s => mtime: %s\n", fpath, fi.ModTime().Format("060102T15:04:05"))
		return nil
	})
	goutil.PanicErr(err)
}
