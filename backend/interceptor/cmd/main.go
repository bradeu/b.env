package main

import (
	"fmt"
	"interceptor/config"
	"interceptor/internal/handlers"
	"interceptor/internal/rabbitmq"
	"interceptor/internal/routes"
	"interceptor/pkg/logger"
	"os"
	"os/signal"
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
	rmq, err := rabbitmq.ConnectRabbitMQ(config.AppConfig.RabbitMQConsumer.URL)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Get the channel
	channel := rmq.GetChannel()

	logger.Info("Channel created successfully")

	// Initialize producer and consumer with exchange and routing key
	producer, err := rabbitmq.NewProducer(
		channel,
		config.AppConfig.RabbitMQProducer.QueueName,
		config.AppConfig.RabbitMQProducer.ExchangeName,
		config.AppConfig.RabbitMQProducer.RoutingKey, // routing key
	)
	if err != nil {
		logger.Fatal("Failed to create producer: %v", err)
	}

	consumer, err := rabbitmq.NewConsumer(
		channel,
		config.AppConfig.RabbitMQConsumer.QueueName,
		config.AppConfig.RabbitMQConsumer.ExchangeName,
		config.AppConfig.RabbitMQConsumer.RoutingKey, // routing key
	)
	if err != nil {
		logger.Fatal("Failed to create consumer: %v", err)
	}

	// // Start consuming messages
	// messages, err := consumer.ConsumeMessages()
	// if err != nil {
	// 	logger.Fatal("Failed to start consuming messages: %v", err)
	// }

	// // Start the consumer in a goroutine
	// go func() {
	// 	if err := handlers.ConsumeMessages(messages); err != nil {
	// 		logger.Error("Consumer error: %v", err)
	// 	}
	// }()

	// logger.Info("Consumer started and listening for messages...")

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
