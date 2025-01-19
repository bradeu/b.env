package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"interceptor/internal/rabbitmq"
	"interceptor/pkg/logger"
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

func ConsumeMessages(messages <-chan amqp.Delivery) (string, error) {
	select {
	case msg := <-messages:
		// Log all message metadata
		logger.Info("Received message with metadata:")
		logger.Info("  Exchange: %s", msg.Exchange)
		logger.Info("  Routing Key: %s", msg.RoutingKey)
		logger.Info("  Content Type: %s", msg.ContentType)
		logger.Info("  Content Encoding: %s", msg.ContentEncoding)
		logger.Info("  Headers:")
		for key, value := range msg.Headers {
			logger.Info("    %s: %v", key, value)
		}
		logger.Info("  Body: %s", string(msg.Body))

		var rawMsg RawMessage
		if err := json.Unmarshal(msg.Body, &rawMsg); err != nil {
			logger.Error("Failed to unmarshal JSON: %v", err)
			msg.Nack(false, true)
			return "", err
		}

		decodedBody, err := base64.StdEncoding.DecodeString(rawMsg.Body)
		if err != nil {
			logger.Error("Failed to decode base64: %v", err)
			msg.Nack(false, true)
			return "", err
		}

		logger.Info("Decoded message: %s", string(decodedBody))
		msg.Ack(false)

		return string(decodedBody), nil

	case <-time.After(5 * time.Second):
		logger.Info("No message available after timeout")
		return "", fmt.Errorf("no message available after timeout")
	}
}
