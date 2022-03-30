package main

import (
	"chat/app"
	"github.com/joho/godotenv"
)

func main() {

	//load .env file from given path
	godotenv.Load(".env")

	app.Start()

	//from := "soen490airbenders@gmail.com"
	//to := "soen390erps@gmail.com"
	//
	//// mailtrap
	////user := "035a001030be3b"
	////password := "5a1e6e53f5f9d8"
	////addr := "smtp.mailtrap.io:2525"
	////host := "smtp.mailtrap.io"
	//
	//// gmail
	//user := "soen490airbenders@gmail.com"
	//password := "airbenders-soen-490"
	//addr := "smtp.gmail.com:587"
	//host := "smtp.gmail.com"
	//
	//msg := "From: soen490airbenders@gmail.com\n" +
	//	"To: soen390erps@gmail.com\n" +
	//	"Subject: Test mail\n\n" +
	//	"Email body"
	//
	//conn, err := net.Dial("tcp", addr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("SUCCESS DIAL SMTP")
	//defer conn.Close()
	//
	//auth := smtp.PlainAuth("", user, password, host)
	//fmt.Println("PLAIN AUTH USED")
	//err = smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Email sent successfully")
}
