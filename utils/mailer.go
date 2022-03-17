package utils

import (
	"net/smtp"
	"os"
)

// Mailer has just one method. To send a simple mail
type Mailer interface {
	SendSimpleMail(to string, body []byte) error
}

// NewSimpleMail is a constructor for Mailer interface. Returns a simpleMail struct
func NewSimpleMail() Mailer {
	return simpleMail{
		from:     os.Getenv("EMAIL_FROM"),
		user:     os.Getenv("USER"),
		password: os.Getenv("PASSWORD"),
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
	}
}

type simpleMail struct {
	from     string
	user     string
	password string
	smtpHost string
	smtpPort string
}

// SendSimpleMail utilizes the golang smtp library to send a simple mail
func (s simpleMail) SendSimpleMail(to string, body []byte) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.smtpHost)
	return smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, body)
}
