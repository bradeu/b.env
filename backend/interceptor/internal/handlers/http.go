package handlers

import (
	"github.com/gofiber/fiber/v2"
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
