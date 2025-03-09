package main

import (
	"fmt"
	"path"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

type GLogConfig137 struct {
	Level            string `yaml:"Level"`
	Pattern          string `yaml:"Pattern"`
	TimeField        string `yaml:"TimeField"`
	TimeFormat       string `yaml:"TimeFormat"`
	Template         string `yaml:"Template"`
	RotateTimeFormat string `yaml:"RotateTimeFormat"`
}

type LogRotateConfig137 struct {
	Filepath   string                `yaml:"filepath"`
	RotateMode rotatefile.RotateMode `yaml:"rotate_mode"`
	RotateTime rotatefile.RotateTime `yaml:"rotate_time"`
	MaxSize    uint64                `yaml:"max_size"`
	BackupNum  uint                  `yaml:"backup_num"`
	BackupTime uint                  `yaml:"backup_time"`
	Compress   bool                  `yaml:"compress"`
	TimeFormat string                `yaml:"time_format"`
	BuffSize   int                   `yaml:"buff_size"`
	BuffMode   string                `yaml:"buff_mode"`
}

type LogConfig137 struct {
	GLogConfig     GLogConfig137      `yaml:"GLogConfig"`
	LogRotate      LogRotateConfig137 `yaml:"LogRotate"`
	ErrorLogRotate LogRotateConfig137 `yaml:"ErrorLogRotate"`
}

func main() {
	slog.DebugMode = true

	logConfig := LogConfig137{
		GLogConfig: GLogConfig137{
			Level:            "debug",
			Pattern:          "development",
			TimeField:        "time",
			TimeFormat:       "2006-01-02 15:04:05.000",
			Template:         "{{datetime}}\t{{level}}\t{{channel}}\t[{{caller}}]\t{{message}}\t{{data}}\t{{extra}}\n",
			RotateTimeFormat: "20060102",
		},
		LogRotate: LogRotateConfig137{
			Filepath:   "testdata/info137c2.log",
			RotateMode: 0,
			RotateTime: 86400,
			MaxSize:    512,
			BackupNum:  3,
			BackupTime: 72,
			Compress:   true,
			TimeFormat: "20060102",
			BuffSize:   512,
			BuffMode:   "line",
		},
		ErrorLogRotate: LogRotateConfig137{
			Filepath:   "testdata/err137c2.log",
			RotateMode: 0,
			RotateTime: 86400,
			MaxSize:    512,
			BackupNum:  3,
			BackupTime: 72,
			Compress:   true,
			TimeFormat: "20060102",
			BuffSize:   512,
			BuffMode:   "line",
		},
	}
	tpl := logConfig.GLogConfig.Template

	// slog.DefaultChannelName = "gookit"
	slog.DefaultTimeFormat = logConfig.GLogConfig.TimeFormat

	slog.Configure(func(l *slog.SugaredLogger) {
		l.Level = slog.TraceLevel
		l.DoNothingOnPanicFatal()
		l.ChannelName = "gookit"
	})
	slog.GetFormatter().(*slog.TextFormatter).SetTemplate(tpl)
	slog.GetFormatter().(*slog.TextFormatter).TimeFormat = slog.DefaultTimeFormat

	rotatefile.DefaultFilenameFn = func(filepath string, rotateNum uint) string {
		suffix := time.Now().Format(logConfig.GLogConfig.RotateTimeFormat)

		// eg: /tmp/error.log => /tmp/error_20250302_01.log
		// 将文件名扩展名取出来, 然后在扩展名中间加入下划线+日期+下划线+序号+扩展名的形式
		ext := path.Ext(filepath)
		filename := filepath[:len(filepath)-len(ext)]

		return filename + fmt.Sprintf("_%s_%02d", suffix, rotateNum) + ext
	}

	h1 := handler.MustRotateFile(logConfig.ErrorLogRotate.Filepath,
		logConfig.ErrorLogRotate.RotateTime,
		// handler.WithFilePerm(os.ModeAppend|os.ModePerm),
		handler.WithLevelMode(slog.LevelModeList),
		handler.WithLogLevels(slog.DangerLevels),
		handler.WithMaxSize(logConfig.ErrorLogRotate.MaxSize),
		handler.WithBackupNum(logConfig.ErrorLogRotate.BackupNum),
		handler.WithBackupTime(logConfig.ErrorLogRotate.BackupTime),
		handler.WithCompress(logConfig.ErrorLogRotate.Compress),
		handler.WithBuffSize(logConfig.ErrorLogRotate.BuffSize),
		handler.WithBuffMode(logConfig.ErrorLogRotate.BuffMode),
		handler.WithRotateMode(logConfig.ErrorLogRotate.RotateMode),
	)
	h1.Formatter().(*slog.TextFormatter).SetTemplate(tpl)

	h2 := handler.MustRotateFile(logConfig.LogRotate.Filepath,
		logConfig.LogRotate.RotateTime,
		// handler.WithFilePerm(os.ModeAppend|os.ModePerm),
		handler.WithLevelMode(slog.LevelModeList),
		handler.WithLogLevels(slog.AllLevels),
		handler.WithMaxSize(logConfig.LogRotate.MaxSize),
		handler.WithBackupNum(logConfig.LogRotate.BackupNum),
		handler.WithBackupTime(logConfig.LogRotate.BackupTime),
		handler.WithCompress(logConfig.LogRotate.Compress),
		handler.WithBuffSize(logConfig.LogRotate.BuffSize),
		handler.WithBuffMode(logConfig.LogRotate.BuffMode),
		handler.WithRotateMode(logConfig.LogRotate.RotateMode),
	)
	h2.Formatter().(*slog.TextFormatter).SetTemplate(tpl)

	slog.PushHandlers(h1, h2)

	// add logs
	for i := 0; i < 20; i++ {
		slog.Infof("hi, this is a example information ... message text. log index=%d", i)
		slog.WithValue("test137", "some value").Warn("测试滚动多个文件，同时设置了清理日志文件")
	}

	slog.MustClose()
	time.Sleep(time.Second * 2)
}
