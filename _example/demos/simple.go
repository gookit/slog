package main

import "github.com/gookit/slog"

// profile run:
//
// go build -gcflags '-m -l' simple.go
func main() {
	// stackIt()
	// _ = stackIt2()
	slogTest()
}

//go:noinline
func stackIt() int {
	y := 2
	return y * 2
}

//go:noinline
func stackIt2() *int {
	y := 2
	res := y * 2
	return &res
}

func slogTest() {
	var msg = "The quick brown fox jumps over the lazy dog"

	slog.Info("rate", "15", "low", 16, "high", 123.2, msg)
	// slog.WithFields(slog.M{
	// 	"omg":    true,
	// 	"number": 122,
	// }).Infof("slog %s", "message message")
}
