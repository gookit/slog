package slog_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
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

// https://github.com/gookit/slog/issues/121
// 当我配置按日期的方式来滚动日志时，当大于 1 天时只能按 1 天来滚动日志。
func TestIssues_121(t *testing.T) {
	seconds := timex.OneDaySec * 7 // 7天
	logFile := "testdata/issue121_7day.log"

	clock := rotatefile.NewMockClock("2024-03-25 08:04:02")
	fh, err := handler.NewTimeRotateFileHandler(
		logFile,
		rotatefile.RotateTime(seconds),
		handler.WithLogLevels(slog.NormalLevels),
		handler.WithBuffSize(128),
		handler.WithBackupNum(20),
		handler.WithTimeClock(clock),
		handler.WithDebugMode, // debug mode
		// handler.WithCompress(log.compress),
		// handler.WithFilePerm(log.filePerm),
	)
	assert.NoError(t, err)

	// create logger with handler and clock.
	l := slog.NewWithHandlers(fh).Config(func(sl *slog.Logger) {
		sl.TimeClock = clock.Now
	})

	// add logs
	for i := 0; i < 50; i++ {
		l.Infof("hi, this is a exmple information ... message text. log index=%d", i)
		clock.Add(24 * timex.Hour)
	}

	l.MustClose()
}

// https://github.com/gookit/slog/issues/137
// 按日期滚动 如果当天时间节点的日志文件已存在 不会append 会直接替换 #137
func TestIssues_137(t *testing.T) {
	logFile := "testdata/issue137_case1.log"
	fsutil.MustSave(logFile, "hello, this is a log file content\n")

	l := slog.NewWithHandlers(handler.MustFileHandler(logFile))

	// add logs
	for i := 0; i < 5; i++ {
		l.Infof("hi, this is a example information ... message text. log index=%d", i)
	}

	l.MustClose()
	// read file content
	content := fsutil.ReadString(logFile)
	assert.StrContains(t, content, "this is a log file content")
	assert.StrContains(t, content, "log index=4")
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

// https://github.com/gookit/slog/issues/144
// slog: failed to handle log, error: write ./logs/info.log: file already closed #144
func TestIssues_144(t *testing.T) {
	defer slog.MustClose()
	slog.Reset()

	// DangerLevels 包含： slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel
	h1 := handler.MustRotateFile("./testdata/logs/error_is144.log", rotatefile.EveryDay,
		handler.WithLogLevels(slog.DangerLevels),
		handler.WithCompress(true),
	)

	// NormalLevels 包含： slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel
	h2 := handler.MustFileHandler("./testdata/logs/info_is144.log", handler.WithLogLevels(slog.NormalLevels))

	// 注册 handler 到 logger(调度器)
	slog.PushHandlers(h1, h2)

	// add logs
	slog.Info("info message text")
	slog.Error("error message text")
}

// https://github.com/gookit/slog/issues/161 自定义level、caller的宽度
func TestIssues_161(t *testing.T) {
	// 这样是全局影响的 - 不推荐
	// slog.LevelNames[slog.WarnLevel] = "WARNI"
	// slog.LevelNames[slog.InfoLevel] = "INFO "
	// slog.LevelNames[slog.NoticeLevel] = "NOTIC"

	l := slog.New()
	l.DoNothingOnPanicFatal()

	h := handler.ConsoleWithMaxLevel(slog.TraceLevel)
	// 通过 SetFormatter 设置格式化 LevelNameLen=5
	h.SetFormatter(slog.TextFormatterWith(slog.LimitLevelNameLen(5)))
	l.AddHandler(h)

	for _, level := range slog.AllLevels {
		l.Logf(level, "a %s test message", level.String())
	}
	assert.NoErr(t, l.LastErr())
}

// https://github.com/gookit/slog/issues/163
func TestIssues_163(t *testing.T) {
	h, e := handler.NewRotateFile("testdata/app_iss163.log", rotatefile.EveryDay)
	assert.NoError(t, e)

	l := slog.NewWithHandlers(h)
	defer l.MustClose()

	l.Debugf("error %+v", e)
	l.Infof("2222")
	// TODO assert.FileExists("testdata/app_iss163.log")
}
