//go:build !windows && !plan9

package handler

import (
	"log/syslog"

	"github.com/gookit/slog"
)

// SysLogOpt for syslog handler
type SysLogOpt struct {
	// Tag syslog tag
	Tag string
	// Priority syslog priority
	Priority syslog.Priority
	// Network syslog network
	Network string
	// Raddr syslog address
	Raddr string
}

// SysLogHandler struct
type SysLogHandler struct {
	slog.LevelWithFormatter
	writer *syslog.Writer
}

// NewSysLogHandler instance
func NewSysLogHandler(priority syslog.Priority, tag string) (*SysLogHandler, error) {
	return NewSysLog(&SysLogOpt{
		Priority: priority,
		Tag:      tag,
	})
}

// NewSysLog handler instance with all custom options.
func NewSysLog(opt *SysLogOpt) (*SysLogHandler, error) {
	slWriter, err := syslog.Dial(opt.Network, opt.Raddr, opt.Priority, opt.Tag)
	if err != nil {
		return nil, err
	}

	h := &SysLogHandler{
		writer: slWriter,
	}

	// init default log level
	h.Level = slog.InfoLevel
	return h, nil
}

// Handle a log record
func (h *SysLogHandler) Handle(record *slog.Record) error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	return h.writer.Info(string(bts))
}

// Close handler
func (h *SysLogHandler) Close() error {
	return h.writer.Close()
}

// Flush handler
func (h *SysLogHandler) Flush() error {
	return nil
}
