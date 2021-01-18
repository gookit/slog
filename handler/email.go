package handler

import (
	"net/smtp"

	"github.com/gookit/slog"
)

// EmailOption struct
type EmailOption struct {
	SmtpHost string // eg "smtp.gmail.com"
	SmtpPort string // eg "587"
	FromAddr string // eg "yourEmail@gmail.com"
	Password string
}

// EmailHandler struct
type EmailHandler struct {
	// NopFlushClose provide empty Flush(), Close() methods
	NopFlushClose
	// LevelWithFormatter support level and formatter
	LevelWithFormatter
	// From the sender email information
	From EmailOption
	// ToAddresses list
	ToAddresses []string
}

// NewEmailHandler instance
func NewEmailHandler(from EmailOption, toAddresses []string) *EmailHandler {
	h := &EmailHandler{
		From: from,
		// to receivers
		ToAddresses: toAddresses,
	}

	// init default log level
	h.Level = slog.InfoLevel
	return h
}

// Handle an log record
func (h *EmailHandler) Handle(r *slog.Record) error {
	msgBytes, err := h.FormatRecord(r)

	var auth = smtp.PlainAuth("", h.From.FromAddr, h.From.Password, h.From.SmtpHost)

	addr := h.From.SmtpHost + ":" + h.From.SmtpPort
	err = smtp.SendMail(addr, auth, h.From.FromAddr, h.ToAddresses, msgBytes)

	return err
}
