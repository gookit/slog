package slog

import (
	"bytes"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/strutil"
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
	case FieldKeyCaller: // eg: "logger_test.go:48,TestLogger_ReportCaller"
		ss := strings.Split(rf.Function, ".")
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line) + ss[len(ss)-1]
	case FieldKeyFile: // eg: "/work/go/gookit/slog/logger_test.go:48"
		return rf.File + ":" + strconv.Itoa(rf.Line)
	case FieldKeyFLine: // eg: "logger_test.go:48"
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line)
	case FieldKeyFcName: // eg: "logger_test.go:48"
		ss := strings.Split(rf.Function, ".")
		return ss[len(ss)-1]
	}

	return rf.File + ":" + strconv.Itoa(rf.Line)
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
func formatArgsWithSpaces(vs []interface{}) (buf []byte) {
	ln := len(vs)
	if ln == 0 {
		return
	}

	if ln == 1 {
		// msg, _ := strutil.AnyToString(vs[0], false)
		msg := stdutil.ToString(vs[0])
		return append(buf, msg...)
	}

	buf = make([]byte, 0, ln*8)
	for i := range vs {
		// str, _ := strutil.AnyToString(vs[i], false)
		str := stdutil.ToString(vs[i])
		if i > 0 { // add space
			buf = append(buf, ' ')
		}
		buf = append(buf, str...)
	}

	return
}

// EncodeToString data to string
func EncodeToString(v interface{}) string {
	if _, ok := v.(map[string]interface{}); ok {
		return mapToString(v.(map[string]interface{}))
	}

	return stdutil.ToString(v)
}

func mapToString(mp map[string]interface{}) string {
	var buf []byte
	buf = append(buf, '{')

	for k, val := range mp {
		buf = append(buf, k...)
		buf = append(buf, ':')

		str, _ := strutil.AnyToString(val, false)
		buf = append(buf, str...)
		buf = append(buf, ',')
	}

	// remove last ','
	buf = append(buf[:len(buf)-1], '}')
	return strutil.Byte2str(buf)
}

func parseTemplateToFields(tplStr string) []string {
	ss := strings.Split(tplStr, "{{")

	vars := make([]string, 0, len(ss)*2)
	for _, s := range ss {
		if len(s) == 0 {
			continue
		}

		fieldAndOther := strings.SplitN(s, "}}", 2)
		if len(fieldAndOther) < 2 {
			vars = append(vars, s)
		} else {
			vars = append(vars, fieldAndOther[0], "}}"+fieldAndOther[1])
		}
	}

	return vars
}
