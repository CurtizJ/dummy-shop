package main

import (
	"net/smtp"
)

type SenderSMTP struct {
	server   string
	login    string
	password string
}

func (sender *SenderSMTP) Send(notification *Notification) error {
	auth := smtp.PlainAuth(
		"",
		sender.login,
		sender.password,
		sender.server,
	)

	msg := "From: " + sender.login + "\n" +
		"To: " + notification.Email + "\n" +
		"Subject: " + notification.Subject + "\n\n" +
		notification.Message

	return smtp.SendMail(
		sender.server+":25",
		auth,
		sender.login,
		[]string{notification.Email},
		[]byte(msg),
	)
}
