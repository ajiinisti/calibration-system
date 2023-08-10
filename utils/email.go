package utils

import (
	"strconv"

	"calibration-system.com/config"
	"gopkg.in/gomail.v2"
)

func SendMail(to []string, subject, message string, cfg config.SMTPConfig) error {

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", cfg.SMTPSenderName)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		port,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}
