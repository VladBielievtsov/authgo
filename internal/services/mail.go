package services

import (
	"authgo/internal/config"
	"log"
	"net/smtp"
	"strings"
)

type MailServices struct {
	cfg  *config.Config
	auth smtp.Auth
}

func NewMailServices(cfg *config.Config) *MailServices {
	if cfg == nil {
		log.Panic("config is nil")
	}

	auth := smtp.PlainAuth(
		"",
		cfg.Mail.Username,
		cfg.Mail.Password,
		cfg.Mail.Host,
	)
	return &MailServices{
		cfg:  cfg,
		auth: auth,
	}
}

func (ms *MailServices) New() smtp.Auth {
	return ms.auth
}

func (ms *MailServices) Send(msg, to string, auth smtp.Auth) error {
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		ms.cfg.Mail.Username,
		[]string{strings.ToLower(to)},
		[]byte(msg),
	)
	return err
}
