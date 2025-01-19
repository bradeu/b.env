package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"solidity/internal/rabbitmq"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

var (
	globalProducer *rabbitmq.Producer
	globalConsumer *rabbitmq.Consumer
)

// HomeHandler responds to the root route
func HomeHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "hello world",
	})
}

// HealthCheckHandler responds to the health check route
func HealthCheckHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Service is healthy",
	})
}

// PublishHandler processes incoming messages
func PublishHandler(c *fiber.Ctx) error {
	// Validate that we actually received a message
	if len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty",
		})
	}

	message := string(c.Body())

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

// InitializeHandlers initializes the global producer and consumer
func InitializeHandlers(producer *rabbitmq.Producer, consumer *rabbitmq.Consumer) {
	globalProducer = producer
	globalConsumer = consumer
}

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

// PublishMessageThroughBroker publishes a message through the broker
func PublishMessageThroughBroker(c *fiber.Ctx) error {
	if len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty",
		})
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        c.Body(),
	}

	if err := globalProducer.PublishMessage(message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// After publishing, consume the message
	messages, err := globalConsumer.ConsumeMessages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to consume messages",
		})
	}

	select {
	case msg := <-messages:
		fmt.Printf("3. Received raw message: %s\n", string(msg.Body))

		// First unmarshal the outer JSON
		var rawMsg RawMessage
		if err := json.Unmarshal(msg.Body, &rawMsg); err != nil {
			fmt.Printf("Failed to unmarshal JSON: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to parse message",
			})
		}

		// Now decode the base64 Body
		decodedBody, err := base64.StdEncoding.DecodeString(rawMsg.Body)
		if err != nil {
			fmt.Printf("4. Decode error: %v\n", err)
			fmt.Printf("4. Failed message: %s\n", rawMsg.Body)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to decode message: %v", err),
			})
		}

		fmt.Printf("5. Decoded message: %s\n", string(decodedBody))

		decoded := DecodedMessage{
			Headers:         msg.Headers,
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			Body:            string(decodedBody),
		}

		msg.Ack(false)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": decoded,
		})

	case <-time.After(10 * time.Second):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No message available",
		})
	}
}

func ConsumeMessageThroughBroker(c *fiber.Ctx) error {
	messages, err := globalConsumer.ConsumeMessages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to consume messages",
		})
	}

	select {
	case msg := <-messages:
		fmt.Printf("3. Received raw message: %s\n", string(msg.Body))

		var rawMsg RawMessage
		if err := json.Unmarshal(msg.Body, &rawMsg); err != nil {
			fmt.Printf("Failed to unmarshal JSON: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to parse message",
			})
		}

		decodedBody, err := base64.StdEncoding.DecodeString(rawMsg.Body)
		if err != nil {
			fmt.Printf("4. Decode error: %v\n", err)
			fmt.Printf("4. Failed message: %s\n", rawMsg.Body)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to decode message: %v", err),
			})
		}

		fmt.Printf("5. Decoded message: %s\n", string(decodedBody))

		decoded := DecodedMessage{
			Headers:         msg.Headers,
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			Body:            string(decodedBody),
		}

		msg.Ack(false)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": decoded,
		})

	case <-time.After(10 * time.Second):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No message available",
		})
	}
}
