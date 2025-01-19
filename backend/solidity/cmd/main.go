package main

import (
	"fmt"
	"os"
	"os/signal"
	"solidity/config"
	"solidity/internal/handlers"
	"solidity/internal/rabbitmq"
	"solidity/internal/routes"
	"solidity/pkg/logger"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables and configuration
	config.LoadEnv()

	// Initialize logger
	if err := logger.InitLogger(config.AppConfig.Logger.FilePath); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logger.Info("Starting application...")

	// 4. Connect to RabbitMQ (which uses logger)
	rmq, err := rabbitmq.ConnectRabbitMQ(config.AppConfig.RabbitMQ.URL)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Get the channel
	channel := rmq.GetChannel()

	// Initialize producer and consumer
	producer, err := rabbitmq.NewProducer(channel, config.AppConfig.RabbitMQ.QueueName)
	if err != nil {
		logger.Fatal("Failed to create producer: %v", err)
	}

	consumer, err := rabbitmq.NewConsumer(channel, config.AppConfig.RabbitMQ.QueueName)
	if err != nil {
		logger.Fatal("Failed to create consumer: %v", err)
	}

	// Initialize handlers
	handlers.InitializeHandlers(producer, consumer)

	// Create a new Fiber app with custom config
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(config.AppConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.AppConfig.Server.WriteTimeout) * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// logger.Error("HTTP Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Register routes
	routes.RegisterRoutes(app)

	// Start server in a goroutine
	go func() {
		serverAddr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
		logger.Info("Server starting on port %s", config.AppConfig.Server.Port)
		if err := app.Listen(serverAddr); err != nil {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}
