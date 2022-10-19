package rabbit

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	err error

	Ch   *amqp.Channel
	Conn *amqp.Connection

	Messages = make(chan string)
	ctx      context.Context
	cancel   context.CancelFunc
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func SetupRabbitMQ() {
	Conn, err = amqp.Dial("amqp://guest:password@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	Ch, err = Conn.Channel()
	failOnError(err, "Failed to open a channel")

	//Declare two queue here ...
	q1, err := Ch.QueueDeclare(
		"taskQueue1", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	q2, err := Ch.QueueDeclare(
		"taskQueue2", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Declare the exchange here use fanout as type for now...
	if err := Ch.ExchangeDeclare(
		"taskExchange", // name of the exchange
		"fanout",       // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		fmt.Printf("Exchange: %s\n", err)
	}

	// Bind Exchange and Queues
	if err := Ch.QueueBind(
		q1.Name,        // name of the queue
		"",             // binding key
		"taskExchange", // source exchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		fmt.Printf("Queue Bind: %s\n", err)
	}

	if err := Ch.QueueBind(
		q2.Name,        // name of the queue
		"",             // binding key
		"taskExchange", // source exchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		fmt.Printf("Queue Bind: %s\n", err)
	}
}

func SendMessage() {
	for msg := range Messages {
		fmt.Println("Got new message")
		body := msg
		err = Ch.PublishWithContext(ctx,
			"taskExchange", // exchange
			"",             // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}

func ReadMessage(queue string) {
	msgs, err := Ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   //args
	)
	if err != nil {
		panic(err)
	}

	for message := range msgs {
		fmt.Printf("Message from %v: %v\n", queue, string(message.Body))
	}
}
