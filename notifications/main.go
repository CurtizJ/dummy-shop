package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	. "github.com/CurtizJ/dummy-shop/lib/config"
	"github.com/streadway/amqp"
)

var sender Sender
var senderMutex sync.Mutex
var config Config

func main() {
	NewConfigFromEnv(&config)

	if config.SenderType == "Email" {
		sender = &SenderSMTP{config.SMTP_SERVER, config.SMTP_LOGIN, config.SMTP_PASSWORD}
	} else {
		sender = &SenderLog{}
	}

	consumer := &AMQPConsumer{}
	amqpURL := os.Getenv("AMQP_ADDR")

	err := consumer.Init(amqpURL, "notifications")
	if err != nil {
		panic("Cannot startup RabbitMQ: " + err.Error())
	}

	err = consumer.Consume(handleMessage)
	if err != nil {
		panic("Cannot start run loop in RabbitMQ: " + err.Error())
	}

	fmt.Println("Waiting for messages at " + amqpURL)
	<-make(chan int)
}

func handleMessage(msg *amqp.Delivery) error {
	var notification Notification
	err := json.Unmarshal(msg.Body, &notification)
	if err != nil {
		return err
	}

	senderMutex.Lock()
	defer senderMutex.Unlock()

	return sender.Send(&notification)
}
