package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func ConsumeMessages() error {
	url := "https://jsonplaceholder.typicode.com/posts/1"

	// Start consuming messages
	messages, err := globalConsumer.ConsumeMessages()
	if err != nil {
		fmt.Printf("Failed to consume messages: %v\n", err)
	}

	// forever := make(chan bool)

	// go func() {
	select {
	case msg := <-messages:
		fmt.Printf("3. Received raw message: %s\n", string(msg.Body))

		var rawMsg RawMessage
		if err := json.Unmarshal(msg.Body, &rawMsg); err != nil {
			fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		}

		decodedBody, err := base64.StdEncoding.DecodeString(rawMsg.Body)
		if err != nil {
			fmt.Printf("4. Decode error: %v\n", err)
			fmt.Printf("4. Failed message: %s\n", rawMsg.Body)
		}

		fmt.Printf("5. Decoded message: %s\n", string(decodedBody))

		decoded := DecodedMessage{
			Headers:         msg.Headers,
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			Body:            string(decodedBody),
		}

		msg.Ack(false)

		logger.Info("Successfully processed message: %v\n", decoded)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Failed to send GET request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v", err)
		}

		fmt.Printf("6. Response body: %s\n", string(body))

		// err = PublishMessage(string(body))
		// if err != nil {
		// 	fmt.Printf("Failed to publish message: %v", err)
		// }

		// <-forever
	}
	// }()

	return nil
}
