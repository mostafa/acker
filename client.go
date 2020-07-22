package main

import (
	"log"

	"github.com/streadway/amqp"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConsumeForever(server string, channel string) {
	if server == "" {
		server = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(server)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	if channel == "" {
		FailOnError(nil, "Channel name is empty")
	}

	queue, err := ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func Produce(server string, channel string, body string, count int) {
	if server == "" {
		server = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(server)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	if channel == "" {
		FailOnError(nil, "Channel name is empty")
	}

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	if count == 0 {
		count = 1
	}

	total := 0
	for i := 0; i < count; i++ {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		if err == nil {
			total += 1
		}
		FailOnError(err, "Failed to publish a message")
	}

	log.Printf(" [x] Published %d messages", total)
}
