package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel   *amqp.Channel
	queueName string
}

func NewProducer(channel *amqp.Channel, queueName string) (*Producer, error) {
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

	return &Producer{
		channel:   channel,
		queueName: queueName,
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
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
