package handler

import (
	"log/syslog"

	"github.com/gookit/slog"
)

type SysLog struct {
	LevelsWithFormatter

	slWriter *syslog.Writer
}

func (h *SysLog) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	return h.slWriter.Info(bts)
}

func NewSysLog(priority syslog.Priority, tag string) *SysLog {
	return &SysLog{
		slWriter: syslog.New(priority, tag),
	}
}
