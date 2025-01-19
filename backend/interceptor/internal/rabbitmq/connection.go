package rabbitmq

import (
	"fmt"
	"interceptor/config"
	"interceptor/pkg/logger"

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
		config.AppConfig.RabbitMQConsumer.QueueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	_, err = r.Channel.QueueDeclare(
		config.AppConfig.RabbitMQProducer.QueueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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
