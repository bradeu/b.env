package handlers

import (
	"encoding/base64"
	"encoding/json"
	"interceptor/internal/rabbitmq"
	"interceptor/pkg/logger"
	"io"
	"net/http"

	"github.com/streadway/amqp"
)

var (
	globalProducer *rabbitmq.Producer
	globalConsumer *rabbitmq.Consumer
)

type DecodedMessage struct {
	Headers         map[string]interface{} `json:"headers"`
	ContentType     string                 `json:"contentType"`
	ContentEncoding string                 `json:"contentEncoding"`
	Body            string                 `json:"body"`
}

type RawMessage struct {
	Headers         map[string]interface{} `json:"Headers"`
	ContentType     string                 `json:"ContentType"`
	ContentEncoding string                 `json:"ContentEncoding"`
	Body            string                 `json:"Body"`
}

// InitializeHandlers initializes the global producer and consumer
func InitializeHandlers(producer *rabbitmq.Producer, consumer *rabbitmq.Consumer) {
	globalProducer = producer
	globalConsumer = consumer
}

func PublishMessage(messages string) error {
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(messages),
	}

	if err := globalProducer.PublishMessage(message); err != nil {
		return err
	}

	return nil
}

func ConsumeMessages(messages <-chan amqp.Delivery) error {
	// Create a channel to keep the consumer running
	forever := make(chan bool)

	go func() {
		for msg := range messages {
			logger.Info("Received raw message: %s", string(msg.Body))

			var rawMsg RawMessage
			if err := json.Unmarshal(msg.Body, &rawMsg); err != nil {
				logger.Error("Failed to unmarshal JSON: %v", err)
				msg.Nack(false, true) // Negative acknowledge and requeue
				continue
			}

			decodedBody, err := base64.StdEncoding.DecodeString(rawMsg.Body)
			if err != nil {
				logger.Error("Failed to decode base64: %v", err)
				msg.Nack(false, true)
				continue
			}

			logger.Info("Decoded message: %s", string(decodedBody))

			// Process decoded message
			resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
			if err != nil {
				logger.Error("Failed to send GET request: %v", err)
				msg.Nack(false, true)
				continue
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read response body: %v", err)
				msg.Nack(false, true)
				continue
			}

			logger.Info("Response body: %s", string(body))
			msg.Ack(false)
		}
	}()

	// Keep the consumer running
	<-forever
	return nil
}
