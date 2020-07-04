package main

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type Notification struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Subject string `json:"subject"`
}

func sendConfirmationLink(email string) error {
	code := strconv.FormatUint(rand.Uint64(), 10)
	if err := confirmations.Set(code, email, time.Second*300).Err(); err != nil {
		return err
	}

	values := url.Values{}
	values.Set("code", code)

	u := url.URL{
		Scheme:   "http",
		Host:     "localhost" + os.Getenv("LISTEN_PORT"),
		Path:     "/auth/confirm",
		RawQuery: values.Encode(),
	}

	message, _ := json.Marshal(&Notification{
		Email:   email,
		Message: "Your confirmation link: " + u.String(),
		Subject: "Email confirmation",
	})

	return notificationsChannel.Publish(
		"",                 // exchange
		notifications.Name, // routing key
		true,               // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
}
