package slog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/strutil"
	"github.com/valyala/bytebufferpool"
)

// const (
// 	defaultMaxCallerDepth  int = 15
// 	defaultKnownSlogFrames int = 4
// )

// var (
// argFmtPool bytebufferpool.Pool
// )

// Stack that attempts to recover the data for all goroutines.
// func getCallStacks(callerSkip int) []byte {
// 	return nil
// }

func buildLowerLevelName() map[Level]string {
	mp := make(map[Level]string, len(LevelNames))
	for level, s := range LevelNames {
		mp[level] = strings.ToLower(s)
	}
	return mp
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

func formatCaller(rf *runtime.Frame, flag uint8) (cs string) {
	switch flag {
	case CallerFlagFull:
		return rf.Function + "," + path.Base(rf.File) + ":" + strconv.FormatInt(int64(rf.Line), 10)
	case CallerFlagFunc:
		return rf.Function
	case CallerFlagFcLine:
		return rf.Function + ":" + strconv.Itoa(rf.Line)
	case CallerFlagPkg:
		i := strings.LastIndex(rf.Function, "/")
		i += strings.IndexByte(rf.Function[i+1:], '.')
		return rf.Function[:i+1]
	case CallerFlagFnlFcn:
		ss := strings.Split(rf.Function, ".")
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line) + "," + ss[len(ss)-1]
	case CallerFlagFnLine:
		return path.Base(rf.File) + ":" + strconv.Itoa(rf.Line)
	case CallerFlagFcName:
		ss := strings.Split(rf.Function, ".")
		return ss[len(ss)-1]
	default: // CallerFlagFpLine
		return rf.File + ":" + strconv.Itoa(rf.Line)
	}
}

// it like Println, will add spaces for each argument
func formatArgsWithSpaces(vs []any) string {
	ln := len(vs)
	if ln == 0 {
		return ""
	}

	if ln == 1 {
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
func EncodeToString(v any) string {
	if mp, ok := v.(map[string]any); ok {
		return mapToString(mp)
	}
	return stdutil.ToString(v)
}

func mapToString(mp map[string]any) string {
	ln := len(mp)
	if ln == 0 {
		return "{}"
	}

	// TODO use bytebufferpool
	buf := make([]byte, 0, ln*8)
	buf = append(buf, '{')

	for k, val := range mp {
		buf = append(buf, k...)
		buf = append(buf, ':')

		str, _ := strutil.AnyToString(val, false)
		buf = append(buf, str...)
		buf = append(buf, ',', ' ')
	}

	// remove last ', '
	buf = append(buf[:len(buf)-2], '}')
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

func printlnStderr(args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}
