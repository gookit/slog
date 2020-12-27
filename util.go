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
	defaultMaxCallerDepth  int = 25
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

func formatCaller(rf *runtime.Frame, field string) (cs string) {
	switch field {
	case FieldKeyCaller: // eg: "logger_test.go:48.TestLogger_ReportCaller"
		ss := strings.Split(rf.Function, ".")
		return fmt.Sprintf("%s:%d.%s", path.Base(rf.File), rf.Line, ss[len(ss)-1])
	case FieldKeyFile: // eg: "/work/go/gookit/slog/logger_test.go:48"
		return fmt.Sprintf("%s:%d", rf.File, rf.Line)
	case FieldKeyFLine: // eg: "logger_test.go:48"
		return fmt.Sprintf("%s:%d", path.Base(rf.File), rf.Line)
	case FieldKeyFcName: // eg: "logger_test.go:48"
		ss := strings.Split(rf.Function, ".")
		return fmt.Sprintf("%s", ss[len(ss)-1])
	}

	return fmt.Sprintf("%s:%d", rf.File, rf.Line)
}

// from glog package
// stacks is a wrapper for runtime.
// Stack that attempts to recover the data for all goroutines.
func getCallStacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit.
	// Start large, though.
	n := 10000
	if all {
		n = 100000
	}

	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		bts := runtime.Stack(trace, all)
		if bts < len(trace) {
			return trace[:bts]
		}
		n *= 2
	}
	return trace
}

// it like Println, will add spaces for each argument
func formatArgsWithSpaces(args []interface{}) (message string) {
	if ln := len(args); ln == 0 {
		message = ""
	} else if ln == 1 {
		message = fmt.Sprint(args[0])
	} else {
		message = fmt.Sprintln(args...)
		// clear last "\n"
		message = message[:len(message)-1]
	}
	return
}
