package rotatefile_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/slog/rotatefile"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain: remove all test files in ./testdata")
	goutil.PanicErr(fsutil.RemoveSub("./testdata", fsutil.ExcludeNames(".keep")))
	m.Run()
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
