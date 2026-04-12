package helpers

import (
	"fmt"
	"net/smtp"

	"github.com/fvrvz/authforest/dto"
	"github.com/fvrvz/gologger"
)

func SendEmail(cfg *dto.SMTP, to string, subject string, body string) error {
	if cfg.Host == "" || cfg.Host == "localhost" {
		gologger.WARN("SMTP not configured. Email to %s with subject '%s' not sent. Body:\n%s", to, subject, body)
		return nil
	}

	from := cfg.From
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, to, subject, body)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		gologger.ERROR("Failed to send email to %s: %v", to, err)
		return err
	}

	return nil
}
