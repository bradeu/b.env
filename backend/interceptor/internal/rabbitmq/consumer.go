package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel      *amqp.Channel
	queueName    string
	exchangeName string
	routingKey   string
}

func NewConsumer(channel *amqp.Channel, queueName, exchangeName, routingKey string) (*Consumer, error) {
	// Declare the exchange
	err := channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare the queue
	_, err = channel.QueueDeclare(
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

	// Bind the queue to the exchange with routing key
	err = channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}, nil
}

func (c *Consumer) ConsumeMessages() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.queueName, // queue name
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // arguments
	)
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

func ConsumeMessages(channelRabbitMQ *amqp.Channel) {
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
		log.Println("Failed to register consumer:", err)
		return
	}

	// Log consumer status
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Listen for messages in a goroutine
	go func() {
		for message := range messages {
			// Process each message
			log.Printf(" > Received message: %s\n", message.Body)

			// Acknowledge if auto-ack is false (for manual acknowledgment)
			// message.Ack(false)
		}
	}()

	// Keep the function running indefinitely
	select {}
}
