package slog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

const (
	defaultMaxCallerDepth int = 25
	defaultKnownSlogFrames int = 4
)

var (
	// qualified package name, cached at first use
	slogPackage string

	// Positions in the call stack when tracing to report the calling method
	minCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// start at the bottom of the stack before the package-name cache is primed
	minCallerDepth = 1
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-slog calling function
func getCaller(maxCallerDepth int) *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maxCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maxCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				slogPackage = getPackageName(funcName)
				break
			}
		}

		minCallerDepth = defaultKnownSlogFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != slogPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func formatCallerToString(rf *runtime.Frame) string {
	// TODO format different string

	// eg: "logger_test.go:48"
	return fmt.Sprintf("%s:%d", path.Base(rf.File), rf.Line)
}

// from glog
// func timeoutFlush(timeout time.Duration) {
// 	done := make(chan bool, 1)
// 	go func() {
// 		FlushAll() // calls logging.lockAndFlushAll()
// 		done <- true
// 	}()
// 	select {
// 	case <-done:
// 	case <-time.After(timeout):
// 		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
// 	}
// }
