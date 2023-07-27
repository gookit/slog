package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/goutil/syncs"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

const pth = "./logs/main.log"

func main() {
	log := slog.New()

	h, err := handler.NewTimeRotateFileHandler(
		pth,
		rotatefile.RotateTime(30),
		handler.WithBuffSize(0),
		handler.WithBackupNum(5),
		handler.WithCompress(true),
		func(c *handler.Config) {
			c.DebugMode = true
		},
	)

	if err != nil {
		panic(err)
	}

	log.AddHandler(h)

	fmt.Println("Start...(can be stop by CTRL+C)", timex.NowDate())
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				log.Info("Log " + time.Now().String())
			}
		}
	}()

	syncs.WaitCloseSignals(func(sig os.Signal) {
		fmt.Println("\nGot signal:", sig)
		fmt.Println("Close logger ...")
		log.MustClose()
	})

	fmt.Println("Exited at", timex.NowDate())
}
