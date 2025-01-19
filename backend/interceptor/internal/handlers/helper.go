package handlers

import (
	"encoding/base64"
	"encoding/json"
	"interceptor/internal/rabbitmq"
	"interceptor/pkg/logger"
	"io"
	"net/http"
	"time"

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
	// Create the message structure
	rawMessage := RawMessage{
		Headers: map[string]interface{}{
			"source":      "interceptor",
			"timestamp":   time.Now().Format(time.RFC3339),
			"messageType": "json",
		},
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Body:            messages,
	}

	// Marshal the raw message
	messageBytes, err := json.Marshal(rawMessage)
	if err != nil {
		logger.Error("Failed to marshal message: %v", err)
		return err
	}

	// Create the AMQP publishing
	message := amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		DeliveryMode:    amqp.Persistent, // Message persistence
		Timestamp:       time.Now(),
		Headers: amqp.Table{
			"source":      "interceptor",
			"messageType": "json",
		},
		Body: messageBytes,
	}

	logger.Info("Attempting to publish message with content: %s", messages)

	// Publish using the global producer
	if err := globalProducer.PublishMessage(message); err != nil {
		logger.Error("Failed to publish message: %v", err)
		return err
	}

	logger.Info("Message published successfully")
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

			PublishMessage(string(body))
		}
	}()

	// Keep the consumer running
	<-forever
	return nil
}
