package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConsumeForever(
	server string, queue string, existing_queue bool, autoack bool, recover bool, currentConsumer bool) {
	if server == "" {
		server = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(server)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	if queue == "" {
		FailOnError(nil, "Queue name is empty")
	}

	if recover {
		ch.Recover(currentConsumer)
	}

	q := amqp.Queue{}

	if existing_queue {
		q, err = ch.QueueDeclarePassive(
			queue, // name
			true,  // durable
			false, // delete when unused (auto_delete)
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		FailOnError(err, "Queue does not exist")
	} else {
		q, err = ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused (auto_delete)
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		FailOnError(err, "Failed to create the queue")
	}

	msgs, err := ch.Consume(
		q.Name,  // queue
		"acker", // consumer
		autoack, // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	total := 0

	go func() {
		for msg := range msgs {
			total += 1
			log.Printf("Received message: #%d, Content: %s", total, msg.Body)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT)
	go func() {
		for sig := range c {
			log.Printf(sig.String())
			log.Printf("Total consumed messages: %d", total)
			os.Exit(-1)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func Produce(server string, queue string, existing_queue bool, body string, count int) {
	if server == "" {
		server = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(server)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	if queue == "" {
		FailOnError(nil, "Queue name is empty")
	}

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q := amqp.Queue{}

	if existing_queue {
		q, err = ch.QueueDeclarePassive(
			queue, // name
			true,  // durable
			false, // delete when unused (auto_delete)
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		FailOnError(err, "Queue does not exist")
	} else {
		q, err = ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused (auto_delete)
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		FailOnError(err, "Failed to declare a queue")
	}

	if count == 0 {
		count = 1
	}

	total := 0
	for i := 0; i < count; i++ {
		err = ch.Publish(
			q.Name, // exchange
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
