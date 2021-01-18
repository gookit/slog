// +build !windows,!nacl,!plan9

package handler

import (
	"log/syslog"

	"github.com/gookit/slog"
)

// SysLogHandler struct
type SysLogHandler struct {
	slWriter *syslog.Writer
	LevelWithFormatter
}

// NewSysLogHandler instance
func NewSysLogHandler(priority syslog.Priority, tag string) (*SysLogHandler, error) {
	slWriter, err := syslog.New(priority, tag)
	if err != nil {
		return nil, err
	}

	h := &SysLogHandler{
		slWriter: slWriter,
	}

	// init default log level
	h.Level = slog.InfoLevel
	return h, nil
}

// Handle an record
func (h *SysLogHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	return h.slWriter.Info(string(bts))
}

func (h *SysLogHandler) Close() error {
	return h.slWriter.Close()
}

func (h *SysLogHandler) Flush() error {
	return nil
}
