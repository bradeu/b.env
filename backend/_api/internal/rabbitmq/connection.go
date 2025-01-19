package rabbitmq

import (
	"api/config"
	"api/pkg/logger"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// ConnectRabbitMQ establishes a connection to RabbitMQ
func ConnectRabbitMQ(url string) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	rmq := &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}

	// Setup queues and exchanges
	if err := rmq.SetupQueuesAndExchanges(); err != nil {
		rmq.Close()
		return nil, err
	}

	logger.Info("Successfully connected to RabbitMQ")
	return rmq, nil
}

// SetupQueuesAndExchanges declares the necessary queues and exchanges
func (r *RabbitMQ) SetupQueuesAndExchanges() error {
	// Declare the queue
	_, err := r.Channel.QueueDeclare(
		config.AppConfig.RabbitMQ.QueueName, // name
		true,                                // durable
		false,                               // delete when unused
		false,                               // exclusive
		false,                               // no-wait
		nil,                                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	return nil
}

// GetChannel returns the channel for use by producer/consumer
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.Channel
}

// Close gracefully closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

// PublishMessage publishes a message to the configured exchange
func (r *RabbitMQ) PublishMessage(message []byte) error {
	return r.Channel.Publish(
		config.AppConfig.RabbitMQ.ExchangeName, // exchange
		"",                                     // routing key
		false,                                  // mandatory
		false,                                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

// ConsumeMessages starts consuming messages from the configured queue
func (r *RabbitMQ) ConsumeMessages(handler func([]byte) error) error {
	messages, err := r.Channel.Consume(
		config.AppConfig.RabbitMQ.QueueName, // queue
		"",                                  // consumer
		false,                               // auto-ack
		false,                               // exclusive
		false,                               // no-local
		false,                               // no-wait
		nil,                                 // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range messages {
			err := handler(msg.Body)
			if err != nil {
				logger.Error("Error processing message: %v", err)
				msg.Nack(false, true) // negative acknowledge and requeue
			} else {
				msg.Ack(false) // acknowledge message
			}
		}
	}()

	return nil
}
