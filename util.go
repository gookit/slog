package slog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/strutil"
	"github.com/valyala/bytebufferpool"
)

// const (
// 	defaultMaxCallerDepth  int = 15
// 	defaultKnownSlogFrames int = 4
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
	lineNum := strconv.FormatInt(int64(rf.Line), 10)
	switch flag {
	case CallerFlagFull:
		return rf.Function + "," + path.Base(rf.File) + ":" + lineNum
	case CallerFlagFunc:
		return rf.Function
	case CallerFlagFcLine:
		return rf.Function + ":" + lineNum
	case CallerFlagPkg:
		i := strings.LastIndex(rf.Function, "/")
		i += strings.IndexByte(rf.Function[i+1:], '.')
		return rf.Function[:i+1]
	case CallerFlagPkgFnl:
		i := strings.LastIndex(rf.Function, "/")
		i += strings.IndexByte(rf.Function[i+1:], '.')
		return rf.Function[:i+1] + "," + path.Base(rf.File) + ":" + lineNum
	case CallerFlagFnlFcn:
		ss := strings.Split(rf.Function, ".")
		return path.Base(rf.File) + ":" + lineNum + "," + ss[len(ss)-1]
	case CallerFlagFnLine:
		return path.Base(rf.File) + ":" + lineNum
	case CallerFlagFcName:
		ss := strings.Split(rf.Function, ".")
		return ss[len(ss)-1]
	default: // CallerFlagFpLine
		return rf.File + ":" + lineNum
	}
}

var msgBufPool bytebufferpool.Pool

// it like Println, will add spaces for each argument
func formatArgsWithSpaces(vs []any) string {
	ln := len(vs)
	if ln == 0 {
		return ""
	}

	if ln == 1 {
		return strutil.SafeString(vs[0]) // perf: Reduce one memory allocation
	}

	// buf = make([]byte, 0, ln*8)
	bb := msgBufPool.Get()
	defer msgBufPool.Put(bb)

	// TIP:
	// `float` to string - will alloc 2 times memory
	// `int <0`, `int > 100` to string -  will alloc 1 times memory
	for i := range vs {
		if i > 0 { // add space
			bb.B = append(bb.B, ' ')
		}
		bb.B = appendAny(bb.B, vs[i])
	}

	return string(bb.B)
	// return byteutil.String(bb.B) // perf: Reduce one memory allocation
}

// TODO replace to byteutil.AppendAny()
func appendAny(dst []byte, v any) []byte {
	if v == nil {
		return append(dst, "<nil>"...)
	}

	switch val := v.(type) {
	case []byte:
		dst = append(dst, val...)
	case string:
		dst = append(dst, val...)
	case int:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int8:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int16:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int32:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int64:
		dst = strconv.AppendInt(dst, val, 10)
	case uint:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint8:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint16:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint32:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint64:
		dst = strconv.AppendUint(dst, val, 10)
	case float32:
		dst = strconv.AppendFloat(dst, float64(val), 'f', -1, 32)
	case float64:
		dst = strconv.AppendFloat(dst, val, 'f', -1, 64)
	case bool:
		dst = strconv.AppendBool(dst, val)
	case time.Time:
		dst = val.AppendFormat(dst, time.RFC3339)
	case time.Duration:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case error:
		dst = append(dst, val.Error()...)
	case fmt.Stringer:
		dst = append(dst, val.String()...)
	default:
		dst = append(dst, fmt.Sprint(v)...)
	}
	return dst
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
