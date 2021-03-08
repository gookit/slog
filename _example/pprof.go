package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// run serve:
// 	go run ./_examples/pprof-cli.go
// see prof on cli:
// 	go tool pprof rux_prof_data.prof
// see prof on web:
// 	go tool pprof -http=:8080 rux_prof_data.prof
func main() {
	logger := slog.NewWithHandlers(
		handler.NewIOWriter(ioutil.Discard, slog.NormalLevels),
	)

	times := 10000
	fmt.Println("start profile, run times:", times)

	ruxProfile := "rux_prof_data.prof"
	f, err := os.Create(ruxProfile)
	if err != nil {
		log.Fatal(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()

	var msg = "The quick brown fox jumps over the lazy dog"
	for i := 0; i < times; i++ {
		logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
	}

	fmt.Println("see prof on web:\n  go tool pprof -http=:8080 rux_prof_data.prof")
}
