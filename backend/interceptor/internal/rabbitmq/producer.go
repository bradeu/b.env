package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel      *amqp.Channel
	queueName    string
	exchangeName string
	routingKey   string
}

func NewProducer(channel *amqp.Channel, queueName, exchangeName, routingKey string) (*Producer, error) {
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

	return &Producer{
		channel:      channel,
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}, nil
}

func (p *Producer) PublishMessage(message interface{}) error {
	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish the message
	return p.channel.Publish(
		p.exchangeName, // exchange
		p.routingKey,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
