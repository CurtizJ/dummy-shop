package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type AMQPConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (consumer *AMQPConsumer) Init(amqpURL, amqpQueue string) error {
	var err error
	consumer.conn, err = amqp.Dial(amqpURL)
	if err != nil {
		return err
	}

	consumer.channel, err = consumer.conn.Channel()
	if err != nil {
		return err
	}

	consumer.queue, err = consumer.channel.QueueDeclare(
		amqpQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return err
}

type MessageCallback func(msg *amqp.Delivery) error

func (consumer *AMQPConsumer) Consume(callback MessageCallback) error {
	messages, err := consumer.channel.Consume(
		consumer.queue.Name, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)

	if err != nil {
		return err
	}

	go consumer.consume(messages, callback)
	return nil
}

func (consumer *AMQPConsumer) consume(messages <-chan amqp.Delivery, callback MessageCallback) {
	for m := range messages {
		go func(msg *amqp.Delivery) {
			err := callback(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s (while handling message: %s)", err.Error(), msg.Body)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}(&m)
	}

	fmt.Println("Finished consuming")
}
