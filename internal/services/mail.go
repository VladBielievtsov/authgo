package services

import (
	"authgo/internal/config"
	"net/smtp"
	"strings"
)

type MailServices struct {
	cfg *config.Config
}

func NewMailServices(cfg *config.Config) *MailServices {
	return &MailServices{cfg: cfg}
}

func (s *MailServices) New() smtp.Auth {
	return smtp.PlainAuth(
		s.cfg.Mail.Identity,
		s.cfg.Mail.Username,
		s.cfg.Mail.Password,
		s.cfg.Mail.Host,
	)
}

func (s *MailServices) Send(msg string, to string, auth smtp.Auth) error {
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		s.cfg.Mail.Username,
		[]string{strings.ToLower(to)},
		[]byte(msg),
	)
	return err
}
