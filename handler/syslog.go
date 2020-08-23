// +build !windows,!nacl,!plan9

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

	return h.slWriter.Info(string(bts))
}

func NewSysLog(priority syslog.Priority, tag string) *SysLog {
	slWriter, err := syslog.New(priority, tag)
	if err != nil {
		panic(err)
	}

	return &SysLog{
		slWriter: slWriter,
	}
}
