package slog

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/strutil"
	"github.com/valyala/bytebufferpool"
)

const (
	defaultMaxCallerDepth  int = 25
	defaultKnownSlogFrames int = 4
)

var (
	// qualified package name, cached at first use. eg: "github.com/gookit/slog"
	slogPackage string

	// Positions in the call stack when tracing to report the calling method
	minCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once

	// argFmtPool bytebufferpool.Pool
)

// Stack that attempts to recover the data for all goroutines.
func getCallStacks(callerSkip int) []byte {
	return nil
}

// getCaller retrieves the name of the first non-slog calling function
func getCaller(callerSkip int) (fr runtime.Frame, ok bool) {
	pcs := make([]uintptr, 1) // alloc 1 times
	num := runtime.Callers(callerSkip, pcs)
	if num < 1 {
		return
	}

	f, _ := runtime.CallersFrames(pcs).Next()
	return f, f.PC != 0
}

func formatCaller(rf *runtime.Frame, field string) (cs string) {
	switch field {
	case FieldKeyCaller:
		// eg: "github.com/gookit/slog_test.TestLogger_ReportCaller,logger_test.go:48"
		return rf.Function + "," + path.Base(rf.File) + ":" + strconv.FormatInt(int64(rf.Line), 10)
	case FieldKeyFLFC:
		ss := strings.Split(rf.Function, ".")
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line) + "," + ss[len(ss)-1]
	case FieldKeyFile: // eg: "/work/go/gookit/slog/logger_test.go:48"
		return rf.File + ":" + strconv.Itoa(rf.Line)
	case FieldKeyFLine:
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line)
	case FieldKeyFcName:
		ss := strings.Split(rf.Function, ".")
		return ss[len(ss)-1]
	}

	return rf.File + ":" + strconv.Itoa(rf.Line)
}

// it like Println, will add spaces for each argument
func formatArgsWithSpaces(vs []interface{}) string {
	ln := len(vs)
	if ln == 0 {
		return ""
	}

	if ln == 1 {
		// msg := stdutil.ToString(vs[0])
		// return strutil.ToBytes(msg) // perf: Reduce one memory allocation
		return stdutil.ToString(vs[0]) // perf: Reduce one memory allocation
	}

	// buf = make([]byte, 0, ln*8)
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	// TIP:
	// `float` to string - will alloc 2 times memory
	// `int <0`, `int > 100` to string -  will alloc 1 times memory
	for i := range vs {
		// str, _ := strutil.AnyToString(vs[i], false)
		str := stdutil.ToString(vs[i])
		if i > 0 { // add space
			// buf = append(buf, ' ')
			bb.B = append(bb.B, ' ')
		}

		// buf = append(buf, str...)
		bb.B = append(bb.B, str...)
	}

	return bb.String()
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
	// TODO use bytebufferpool
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
