package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
	Logger   LoggerConfig
}

// ServerConfig holds all HTTP server related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// RabbitMQConfig holds all RabbitMQ related configuration
type RabbitMQConfig struct {
	URL          string
	QueueName    string
	ExchangeName string
}

// LoggerConfig holds all logger related configuration
type LoggerConfig struct {
	FilePath string
	MinLevel string
}

var AppConfig Config

// LoadEnv loads the environment variables from .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	// Initialize the global config
	AppConfig = Config{
		Server: ServerConfig{
			Port:         GetEnv("SERVER_PORT", "3000"),
			ReadTimeout:  GetEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: GetEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
		},
		RabbitMQ: RabbitMQConfig{
			URL:          GetEnv("AMQP_SERVER_URL", "amqp://guest:guest@localhost:5672/"),
			QueueName:    GetEnv("AMQP_QUEUE_NAME", "default_queue"),
			ExchangeName: GetEnv("AMQP_EXCHANGE_NAME", "default_exchange"),
		},
		Logger: LoggerConfig{
			FilePath: GetEnv("LOG_FILE_PATH", filepath.Join("logs", "app.log")),
			MinLevel: GetEnv("LOG_MIN_LEVEL", "DEBUG"),
		},
	}
}

// GetEnv retrieves an environment variable with a fallback value
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvAsInt retrieves an environment variable as integer with a fallback value
func GetEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
