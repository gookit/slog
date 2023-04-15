package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Obj struct {
	a int
	b int64
	c string
	d bool
}

var (
	str1 = "str1"
	str2 = "str222222222222"
	int1 = 1
	int2 = 2
	obj  = Obj{1, 2, "3", true}
)

func TestZapSugar(t *testing.T) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./zap-sugar.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	sugar := logger.Sugar()
	sugar.Info("message is msg")

	count := 100000
	start := time.Now().UnixNano()
	for n := count; n > 0; n-- {
		sugar.Info("message is msg")
	}
	end := time.Now().UnixNano()
	fmt.Printf("\n zap sugar no format\n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)

	start = time.Now().UnixNano()
	for n := count; n > 0; n-- {
		sugar.Infof("message is %d %d %s %s %#v", int1, int2, str1, str2, obj)
	}
	end = time.Now().UnixNano()
	fmt.Printf("\n zap sugar format\n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)
	sugar.Sync()
}

func TestZapLog(t *testing.T) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	count := 100000
	start := time.Now().UnixNano()
	for n := count; n > 0; n-- {
		logger.Info("message is msg")
	}
	end := time.Now().UnixNano()
	fmt.Printf("\n zap no format\n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)

	start = time.Now().UnixNano()
	for n := count; n > 0; n-- {
		logger.Info("failed to fetch URL",
			// Structured context as strongly typed Field values.
			zap.Int("int1", int1),
			zap.Int("int2", int2),
			zap.String("str", str1),
			zap.String("str2", str2),
			zap.Any("backoff", obj),
		)
	}
	end = time.Now().UnixNano()
	fmt.Printf("\n zap format\n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)
	logger.Sync()
}

func TestSlog(t *testing.T) {
	h1, err := handler.NewEmptyConfig(
		handler.WithLogfile("./slog-info.log"),    // 路径
		handler.WithRotateTime(handler.EveryHour), // 日志分割间隔
		handler.WithLogLevels(slog.AllLevels),     // 日志level
		handler.WithBuffSize(4*1024*1024),         // buffer大小
		handler.WithCompress(true),                // 是否压缩旧日志 zip
		handler.WithBackupNum(24*3),               // 保留旧日志数量
		handler.WithBuffMode(handler.BuffModeBite),
		// handler.WithRenameFunc(),                    //RenameFunc build filename for rotate file
	).CreateHandler()
	if err != nil {
		fmt.Printf("Create slog handler err: %#v", err)
		return
	}

	f := slog.AsTextFormatter(h1.Formatter())
	myTplt := "[{{datetime}}] [{{level}}] [{{caller}}] {{message}}\n"
	f.SetTemplate(myTplt)
	logs := slog.NewWithHandlers(h1)

	count := 100000
	start := time.Now().UnixNano()
	for i := 0; i < count; i++ {
		logs.Info("message is msg")
	}
	end := time.Now().UnixNano()
	fmt.Printf("\n slog no format \n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)

	start = time.Now().UnixNano()
	for n := count; n > 0; n-- {
		logs.Infof("message is %d %d %s %s %#v", int1, int2, str1, str2, obj)
	}
	end = time.Now().UnixNano()
	fmt.Printf("\n slog format \n total cost %d ns\n  avg  cost %d ns \n count %d \n", end-start, (end-start)/int64(count), count)
	logs.MustFlush()
}
