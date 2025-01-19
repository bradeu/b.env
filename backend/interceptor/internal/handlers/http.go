package handlers

import (
	"fmt"
	"interceptor/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

// var (
// 	globalProducer *rabbitmq.Producer
// 	globalConsumer *rabbitmq.Consumer
// )

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

// // InitializeHandlers initializes the global producer and consumer
// func InitializeHandlers(producer *rabbitmq.Producer, consumer *rabbitmq.Consumer) {
// 	globalProducer = producer
// 	globalConsumer = consumer
// }

// PublishMessageThroughBroker publishes a message through the broker
func PublishMessageThroughBroker(c *fiber.Ctx) error {
	if len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty",
		})
	}

	messages := string(c.Body())

	// Create a channel to receive the result
	resultChan := make(chan struct {
		msg string
		err error
	})

	// Publish and consume in a goroutine
	go func() {
		// Publish the message
		if err := PublishMessage(messages); err != nil {
			resultChan <- struct {
				msg string
				err error
			}{"", fmt.Errorf("failed to publish message: %v", err)}
			return
		}

		// Get the message channel
		msgChan, err := globalConsumer.ConsumeMessages()
		if err != nil {
			resultChan <- struct {
				msg string
				err error
			}{"", fmt.Errorf("failed to start consuming messages: %v", err)}
			return
		}

		// Try to consume the message and get the response
		consumedMsg, err := ConsumeMessages(msgChan)
		resultChan <- struct {
			msg string
			err error
		}{consumedMsg, err}
	}()

	// Wait for the result with a timeout
	select {
	case result := <-resultChan:
		if result.err != nil {
			logger.Error("Operation failed: %v", result.err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": result.err.Error(),
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"status":   "success",
			"sent":     messages,
			"received": result.msg,
		})

	case <-time.After(10 * time.Second): // Adjust timeout as needed
		return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
			"status":  "error",
			"message": "Operation timed out",
		})
	}
}
