package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/rabbitmq"
	"api/internal/routes"
	"api/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	// Start consuming messages
	messages, err := consumer.ConsumeMessages()
	if err != nil {
		fmt.Printf("Failed to consume messages: %v\n", err)
	}

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

	case <-time.After(10 * time.Second):
		fmt.Println("No message available")
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
