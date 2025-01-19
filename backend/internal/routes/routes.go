package routes

import (
	"backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// API routes group
	api := app.Group("/api")

	// Public routes
	app.Get("/", handlers.HomeHandler)
	app.Get("/health", handlers.HealthCheckHandler)

	// API endpoints
	api.Post("/publish", handlers.PublishHandler)

	// RabbitMQ endpoints
	api.Post("/publishbroker", handlers.PublishMessageThroughBroker)
}
