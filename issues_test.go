package slog_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// https://github.com/gookit/slog/issues/27
func TestIssues_27(t *testing.T) {
	defer slog.Reset()

	count := 0
	for {
		if count >= 6 {
			break
		}
		slog.Infof("info log %d", count)
		time.Sleep(time.Second)
		count++
	}
}

// https://github.com/gookit/slog/issues/31
func TestIssues_31(t *testing.T) {
	defer slog.Reset()
	defer slog.MustClose()

	// slog.DangerLevels equals slog.Levels{slog.PanicLevel, slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel}
	h1 := handler.MustFileHandler("testdata/error_issue31.log", handler.WithLogLevels(slog.DangerLevels))

	infoLevels := slog.Levels{slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
	h2 := handler.MustFileHandler("testdata/info_issue31.log", handler.WithLogLevels(infoLevels))

	slog.PushHandler(h1)
	slog.PushHandlers(h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}

// https://github.com/gookit/slog/issues/52
func TestIssues_52(t *testing.T) {
	testTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}"
	slog.SetLogLevel(slog.ErrorLevel)
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(testTemplate)

	slog.Error("Error message")
	slog.Reset()

	fmt.Println()
	// dump.P(slog.GetFormatter())
}

// https://github.com/gookit/slog/issues/75
func TestIssues_75(t *testing.T) {
	slog.Error("Error message 1")

	// set max level
	slog.SetLogLevel(slog.Level(0))
	// slog.SetLogLevel(slog.PanicLevel)
	slog.Error("Error message 2")
	slog.Reset()
	// dump.P(slog.GetFormatter())
}

// https://github.com/gookit/slog/issues/105
func TestIssues_105(t *testing.T) {
	t.Run("simple write", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			slog.Error("simple error log", i)
			time.Sleep(time.Millisecond * 100)
		}
	})

	// test concurrent write
	t.Run("concurrent write", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				slog.Error("concurrent error log", i)
				time.Sleep(time.Millisecond * 100)
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}

// https://github.com/gookit/slog/issues/108
func TestIssues_108(t *testing.T) {
	buf1 := byteutil.NewBuffer()
	root := slog.NewWithName("root", func(l *slog.Logger) {
		l.ChannelName = l.Name()
		l.AddHandler(handler.NewSimple(buf1, slog.InfoLevel))
	})
	root.Info("root info message")
	root.Warn("root warn message")

	str := buf1.ResetGet()
	fmt.Println(str)
	assert.StrContains(t, str, "[root] [INFO")
	assert.StrContains(t, str, "[root] [WARN")

	buf2 := byteutil.NewBuffer()
	probe := slog.NewWithName("probe", func(l *slog.Logger) {
		l.ChannelName = l.Name()
		l.AddHandler(handler.NewSimple(buf2, slog.InfoLevel))
	})
	probe.Info("probe info message")
	probe.Warn("probe warn message")

	str = buf2.ResetGet()
	fmt.Println(str)
	assert.StrContains(t, str, "[probe] [INFO")
	assert.StrContains(t, str, "[probe] [WARN")
}

// https://github.com/gookit/slog/issues/139
// 自定义模板报 invalid memory address or nil pointer dereference #139
func TestIssues_139(t *testing.T) {
	myTemplate := "[{{datetime}}] [{{requestid}}] [{{level}}] {{message}}\n"
	textFormatter := &slog.TextFormatter{TimeFormat: "2006-01-02 15:04:05.000"}
	textFormatter.SetTemplate(myTemplate)
	// use func create
	// textFormatter := slog.NewTextFormatter(myTemplate).Configure(func(f *slog.TextFormatter) {
	// 	f.TimeFormat = "2006-01-02 15:04:05.000"
	// })
	h1 := handler.NewConsoleHandler(slog.AllLevels)
	h1.SetFormatter(textFormatter)

	L := slog.New()
	L.AddHandlers(h1)
	// add processor <====
	// L.AddProcessor(slog.ProcessorFunc(func(r *slog.Record) {
	// 	r.Fields["requestid"] = r.Ctx.Value("requestid")
	// }))
	L.AddProcessor(slog.AppendCtxKeys("requestid"))

	ctx := context.WithValue(context.Background(), "requestid", "111111")
	L.WithCtx(ctx).Info("test")
}
