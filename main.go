package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {

	//load .env file from given path
	//godotenv.Load(".env")

	//app.Start()

	from := "john.doe@example.com"
	user := "035a001030be3b"
	password := "5a1e6e53f5f9d8"

	to := []string{
		"roger.roe@example.com",
	}

	addr := "smtp.mailtrap.io:2525"
	host := "smtp.mailtrap.io"

	msg := []byte("From: john.doe@example.com\r\n" +
		"To: roger.roe@example.com\r\n" +
		"Subject: Test mail\r\n\r\n" +
		"Email body\r\n")

	auth := smtp.PlainAuth("", user, password, host)

	err := smtp.SendMail(addr, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
