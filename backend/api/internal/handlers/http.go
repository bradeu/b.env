package handlers

import (
	"api/internal/rabbitmq"
	"api/pkg/logger"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	globalProducer *rabbitmq.Producer
	globalConsumer *rabbitmq.Consumer
)

// InitializeHandlers initializes the global producer and consumer
func InitializeHandlers(producer *rabbitmq.Producer, consumer *rabbitmq.Consumer) {
	globalProducer = producer
	globalConsumer = consumer
}

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

// PublishMessageThroughBroker publishes a message through the broker
func PublishMessageThroughBroker(c *fiber.Ctx) error {
	if len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty",
		})
	}

	// Try to parse as JSON first
	var jsonData map[string]interface{}
	err := json.Unmarshal(c.Body(), &jsonData)
	if err != nil {
		// If parsing fails, treat as raw message
		jsonData = map[string]interface{}{
			"content": string(c.Body()),
		}
	}

	// Add metadata
	message := map[string]interface{}{
		"data":      jsonData,
		"timestamp": time.Now(),
	}

	// Publish message
	if err := globalProducer.PublishMessage(message); err != nil {
		logger.Error("Failed to publish message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to publish message",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":  "success",
		"message": "Message published successfully",
	})
}

// PublishBatchThroughBroker publishes a batch of messages through the broker
func PublishBatchThroughBroker(c *fiber.Ctx) error {
	// Parse batch messages
	var messages []interface{}
	if err := json.Unmarshal(c.Body(), &messages); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid batch format. Expected array of messages",
		})
	}

	// Process each message
	results := make([]fiber.Map, len(messages))
	for i, msg := range messages {
		// Add metadata
		enrichedMsg := map[string]interface{}{
			"data":      msg,
			"timestamp": time.Now(),
			"batch_id":  i,
		}

		// Publish message
		err := globalProducer.PublishMessage(enrichedMsg)
		if err != nil {
			logger.Error("Failed to publish batch message %d: %v", i, err)
			results[i] = fiber.Map{
				"status":  "error",
				"message": "Failed to publish message",
				"index":   i,
			}
		} else {
			results[i] = fiber.Map{
				"status":  "success",
				"message": "Message published successfully",
				"index":   i,
			}
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":  "success",
		"results": results,
	})
}

// GetQueueStatus returns information about the message queue
func GetQueueStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"queue": fiber.Map{
			"name":    globalProducer.GetQueueName(),
			"active":  true,
			"running": true,
		},
	})
}
