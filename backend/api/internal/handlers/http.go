package handlers

import (
	"api/internal/rabbitmq"
	"fmt"

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

	messages, err := globalConsumer.ConsumeMessages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to consume messages",
		})
	}

	go func() {
		for msg := range messages {
			fmt.Printf("Received message: %s\n", string(msg.Body))
			msg.Ack(false)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "message sent to broker",
	})
}
