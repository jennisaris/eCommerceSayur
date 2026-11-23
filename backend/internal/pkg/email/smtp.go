package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPSender struct {
	config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{config: config}
}

func (s *SMTPSender) IsConfigured() bool {
	return strings.TrimSpace(s.config.Host) != "" &&
		strings.TrimSpace(s.config.Username) != "" &&
		strings.TrimSpace(s.config.Password) != "" &&
		strings.TrimSpace(s.config.From) != ""
}

func sanitizeHeaderValue(value string) string {
	cleaned := strings.ReplaceAll(value, "\r", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	return strings.TrimSpace(cleaned)
}

func (s *SMTPSender) SendVerificationEmail(ctx context.Context, recipient, recipientName, verificationURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	safeFrom := sanitizeHeaderValue(s.config.From)
	safeRecipient := sanitizeHeaderValue(recipient)
	safeName := sanitizeHeaderValue(recipientName)
	safeURL := sanitizeHeaderValue(verificationURL)

	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Verify your BookVerse account\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nHello %s,\r\n\r\nVerify your account using this link:\r\n%s\r\n\r\nThis link expires in 24 hours.\r\n",
		safeFrom, safeRecipient, safeName, safeURL,
	))
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	return smtp.SendMail(s.config.Host+":"+s.config.Port, auth, safeFrom, []string{safeRecipient}, message)
}

func (s *SMTPSender) SendPasswordResetEmail(ctx context.Context, recipient, recipientName, resetURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	safeFrom := sanitizeHeaderValue(s.config.From)
	safeRecipient := sanitizeHeaderValue(recipient)
	safeName := sanitizeHeaderValue(recipientName)
	safeURL := sanitizeHeaderValue(resetURL)

	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Reset your BookVerse password\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nHello %s,\r\n\r\nWe received a request to reset your BookVerse password. Use the link below to set a new password:\r\n%s\r\n\r\nThis link expires in 1 hour. If you did not request this, you can safely ignore this email.\r\n",
		safeFrom, safeRecipient, safeName, safeURL,
	))
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	return smtp.SendMail(s.config.Host+":"+s.config.Port, auth, safeFrom, []string{safeRecipient}, message)
}
