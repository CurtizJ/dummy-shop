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
}

func sendConfirmationLink(email string) error {
	code := rand.Uint64()
	if err := confirmations.Set(email, code, time.Second*60).Err(); err != nil {
		return err
	}

	values := url.Values{}
	values.Set("email", email)
	values.Set("code", strconv.FormatUint(code, 10))

	u := url.URL{
		Scheme:   "http",
		Host:     "localhost" + os.Getenv("LISTEN_PORT"),
		Path:     "/auth/confirm",
		RawQuery: values.Encode(),
	}

	message, _ := json.Marshal(&Notification{
		Email:   email,
		Message: "Your confirmation link: " + u.String(),
	})

	return notificationsChannel.Publish(
		"",                 // exchange
		notifications.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
