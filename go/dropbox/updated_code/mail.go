package main

import (
	"log"
	"net/smtp"
)

func main() {
	send("Yay!")
}

func send(body string) {
	from := "boxdrop162@gmail.com"
	pass := "notmyactualpassword"
	to := "jaredjtc@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: This worked!\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}
