package slog

import (
	"fmt"
	"io"
	"testing"

	"github.com/gookit/goutil/dump"
)

func TestLogger_newRecord_AllocTimes(_ *testing.T) {
	l := Std()
	l.Output = io.Discard
	defer l.Reset()

	// output: 0 times
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(100, func() {
		// logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
		r := l.newRecord()
		// do something...
		l.releaseRecord(r)
	})))
}

func Test_AllocTimes_formatArgsWithSpaces_oneElem(_ *testing.T) {
	// string Alloc Times: 0
	fmt.Println("string Alloc Times:", int(testing.AllocsPerRun(10, func() {
		// logger.Info("rate", "15", "low", 16, "high", 123.2, msg)
		formatArgsWithSpaces([]any{"msg"})
	})))

	// int Alloc Times: 1
	fmt.Println("int Alloc Times:", int(testing.AllocsPerRun(10, func() {
		formatArgsWithSpaces([]any{2343})
	})))

	// float Alloc Times: 2
	fmt.Println("float Alloc Times:", int(testing.AllocsPerRun(10, func() {
		formatArgsWithSpaces([]any{123.2})
	})))
}

func Test_AllocTimes_formatArgsWithSpaces_manyElem(_ *testing.T) {
	// Alloc Times: 1
	// TIP:
	// `float` will alloc 2 times memory
	// `int <0`, `int > 100` will alloc 1 times memory
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(50, func() {
		formatArgsWithSpaces([]any{
			"rate", -23, true, 106, "high", 123.2,
		})
	})))
}

func Test_AllocTimes_stringsPool(_ *testing.T) {
	l := Std()
	l.Output = io.Discard
	l.LowerLevelName = true
	defer l.Reset()

	var ln, cp int
	// output: 0 times
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(100, func() {
		// logger.Info("rate", "15", "low", 16, "high", 123.2, msg)

		// oldnew := stringsPool.Get().([]string)
		// defer stringsPool.Put(oldnew)

		oldnew := make([]string, 0, len(map[string]string{"a": "b"})*2+1)

		oldnew = append(oldnew, "a")
		oldnew = append(oldnew, "b")
		oldnew = append(oldnew, "c")
		// oldnew = append(oldnew, "d")

		ln = len(oldnew)
		cp = cap(oldnew)
	})))

	dump.P(ln, cp)
}

func TestLogger_Info_oneElem_AllocTimes(_ *testing.T) {
	l := Std()
	// l.Output = io.Discard
	l.ReportCaller = false
	l.LowerLevelName = true
	// 启用 color 会导致多次(10次左右)内存分配
	l.Formatter.(*TextFormatter).EnableColor = false

	defer l.Reset()

	// output: 2 times
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(5, func() {
		// l.Info("rate", "15", "low", 16, "high", 123.2, "msg")
		l.Info("msg")
	})))
}

func TestLogger_Info_moreElem_AllocTimes(_ *testing.T) {
	l := NewStdLogger()
	// l.Output = io.Discard
	l.ReportCaller = false
	l.LowerLevelName = true
	// 启用 color 会导致多次(10次左右)内存分配
	l.Formatter.(*TextFormatter).EnableColor = false

	defer l.Reset()

	// output: 5 times
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(5, func() {
		l.Info("rate", "15", "low", 16, "high", 123.2, "msg")
	})))

	// output: 5 times
	fmt.Println("Alloc Times:", int(testing.AllocsPerRun(5, func() {
		l.Info("rate", "15", "low", 16, "high")
		// l.Info("msg")
	})))
}
