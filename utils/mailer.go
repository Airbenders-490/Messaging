package utils

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"
)

// Mailer has just one method. To send a simple mail
type Mailer interface {
	SendSimpleMail(to string, body []byte) error
}

// NewSimpleMail is a constructor for Mailer interface. Returns a simpleMail struct
func NewSimpleMail() Mailer {
	mailer := simpleMail{
		from:     os.Getenv("EMAIL_FROM"),
		user:     os.Getenv("USER"),
		password: os.Getenv("PASSWORD"),
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
	}
	fmt.Sprintf("SIMPLE MAIL CONFIG: %s %s %s %s %s", mailer.from, mailer.user, mailer.password, mailer.password, mailer.smtpHost, mailer.smtpHost)
	return mailer
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
	fmt.Sprintf("EMAIL CREDS: %s , %s , %s , %s , %s\n", s.from, s.user, s.password, s.smtpHost, s.smtpPort)
	// PlainAuth will only send the credentials if the connection is using TLS or is
	// connected to localhost.
	// Otherwise authentication will fail with an error, without sending the credentials.
	auth := smtp.PlainAuth("", s.user, s.password, s.smtpHost)
	fmt.Sprintf("COMPLETED AUTH SETUP %s\n", auth)

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, body)
	if err != nil {
		fmt.Sprintf("FAILED TO SEND MAIL WITH PLAINAUTH\n %s", err)

		auth = smtp.CRAMMD5Auth(s.user, s.password)
		err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, body)
		if err != nil {
			fmt.Sprintf("FAILED TO SEND MAIL USING CRAMMD5AUTH\n %s", err)

			err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, nil, s.from, []string{to}, body)

			if err != nil {
				fmt.Sprintf("FAILED TO SEND MAIL WITHOUT AUTH\n %s", err)


				_, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort))
				if err != nil {
					fmt.Sprintf("COULD NOT DIAL IN TO SMTP ADDRESS\n %s", err)

					conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort), 10*time.Second)
					if err != nil {
						fmt.Sprintf("COULD NOT DIALTIMEOUT IN TO SMTP ADDRESS\n %s", err)
					}
					fmt.Sprintf("SUCCESS DIALTIMEOUT TO SMTP ADDRESS\n %s", err)
					// Connect to the SMTP server
					c, err := smtp.NewClient(conn, s.smtpHost)
					if err != nil {
						fmt.Sprintf("FAILED TO CREATE NEW STMP CLIENT \n %s", err)
					}
					defer c.Quit()
				} else {
					fmt.Println("SUCCESSFULLY DIALLED INTO SMTP ADDRESS")
				}
			} else {
				fmt.Sprintf("SUCCESS SEND MAIL WITHOUT AUTH\n")
			}


		}else {
			fmt.Sprintf("SUCCESS SEND MAIL USING CRAMMD5AUTH\n")
		}
	} else {
		fmt.Println("SUCCESS SENT MAIL USING PLAINAUTH")
	}
	return err
}
