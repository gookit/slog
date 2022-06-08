package handler

import (
	"net/smtp"
	"strconv"

	"github.com/gookit/slog"
)

// EmailOption struct
type EmailOption struct {
	SMTPHost string `json:"smtp_host"` // eg "smtp.gmail.com"
	SMTPPort int    `json:"smtp_port"` // eg 587
	FromAddr string `json:"from_addr"` // eg "yourEmail@gmail.com"
	Password string `json:"password"`
}

// EmailHandler struct
type EmailHandler struct {
	NopFlushClose
	slog.LevelWithFormatter
	// From the sender email information
	From EmailOption
	// ToAddresses email list
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

// Handle a log record
func (h *EmailHandler) Handle(r *slog.Record) error {
	msgBytes, err := h.Format(r)
	if err != nil {
		return err
	}

	var auth = smtp.PlainAuth("", h.From.FromAddr, h.From.Password, h.From.SMTPHost)

	addr := h.From.SMTPHost + ":" + strconv.Itoa(h.From.SMTPPort)

	return smtp.SendMail(addr, auth, h.From.FromAddr, h.ToAddresses, msgBytes)
}
