package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel   *amqp.Channel
	queueName string
}

func NewConsumer(channel *amqp.Channel, queueName string) (*Consumer, error) {
	// Declare the queue
	_, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *Consumer) ConsumeMessages() (<-chan amqp.Delivery, error) {
	// Use the existing channel instead of creating new one
	messages, err := c.channel.Consume(
		c.queueName, // queue name
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Consumer) Consume(handler func([]byte) error) error {
	messages, err := c.ConsumeMessages()
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for msg := range messages {
			if err := handler(msg.Body); err != nil {
				// Handle error (log it, retry, etc.)
				// For now, just print it
				println("Error processing message:", err.Error())
			}
		}
	}()

	<-forever
	return nil
}

func ConsumeMessages() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		"QueueService1", // queue name
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever
}
