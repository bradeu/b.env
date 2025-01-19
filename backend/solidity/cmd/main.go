package main

import (
	"fmt"
	"os"
	"os/signal"
	"solidity/config"
	"solidity/internal/handlers"
	"solidity/internal/rabbitmq"
	"solidity/pkg/logger"
	"syscall"
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

	logger.Info("Channel created successfully")
	logger.Info("Exchange: %s", config.AppConfig.RabbitMQ.ExchangeName)
	logger.Info("Queue: %s", config.AppConfig.RabbitMQ.QueueName)
	logger.Info("Routing Key: %s", config.AppConfig.RabbitMQ.RoutingKey)

	// Initialize producer and consumer with exchange and routing key
	producer, err := rabbitmq.NewProducer(
		channel,
		config.AppConfig.RabbitMQ.QueueName,
		config.AppConfig.RabbitMQ.ExchangeName,
		config.AppConfig.RabbitMQ.RoutingKey, // routing key
	)
	if err != nil {
		logger.Fatal("Failed to create producer: %v", err)
	}

	consumer, err := rabbitmq.NewConsumer(
		channel,
		config.AppConfig.RabbitMQ.QueueName,
		config.AppConfig.RabbitMQ.ExchangeName,
		config.AppConfig.RabbitMQ.RoutingKey, // routing key
	)
	if err != nil {
		logger.Fatal("Failed to create consumer: %v", err)
	}

	// Start consuming messages
	messages, err := consumer.ConsumeMessages()
	if err != nil {
		logger.Fatal("Failed to start consuming messages: %v", err)
	}

	// Start the consumer in a goroutine
	go func() {
		if err := handlers.ConsumeMessages(messages); err != nil {
			logger.Error("Consumer error: %v", err)
		}
	}()

	logger.Info("Consumer started and listening for messages...")

	// Initialize handlers
	handlers.InitializeHandlers(producer, consumer)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}
