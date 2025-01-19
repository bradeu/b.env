package services

import (
	"api/pkg/logger"
	"encoding/json"
	"fmt"
	"time"
)

// MessageProcessor handles the processing of messages
type MessageProcessor struct {
	// Add any dependencies here (e.g., database, external services)
}

// Message represents the structure of messages we expect to process
type Message struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMessageProcessor creates a new instance of MessageProcessor
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{}
}

// ProcessMessage handles the processing of a single message
func (p *MessageProcessor) ProcessMessage(messageBytes []byte) error {
	// Parse the message
	var message Message
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Log the received message
	logger.Info("Processing message: ID=%s, Type=%s", message.ID, message.Type)

	// Process the message based on its type
	switch message.Type {
	case "task":
		return p.processTask(message)
	case "notification":
		return p.processNotification(message)
	default:
		return fmt.Errorf("unknown message type: %s", message.Type)
	}
}

// processTask handles task-type messages
func (p *MessageProcessor) processTask(message Message) error {
	// Extract task data
	taskData, ok := message.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid task data format")
	}

	// Log task processing
	logger.Info("Processing task: %v", taskData)

	// TODO: Implement task processing logic
	// This could include:
	// - Saving to database
	// - Calling external services
	// - Updating application state
	// - Triggering other processes

	return nil
}

// processNotification handles notification-type messages
func (p *MessageProcessor) processNotification(message Message) error {
	// Extract notification data
	notifData, ok := message.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid notification data format")
	}

	// Log notification processing
	logger.Info("Processing notification: %v", notifData)

	// TODO: Implement notification processing logic
	// This could include:
	// - Sending emails
	// - Pushing to websockets
	// - Triggering alerts
	// - Updating notification status

	return nil
}

// ValidateMessage checks if a message has all required fields
func (p *MessageProcessor) ValidateMessage(message Message) error {
	if message.ID == "" {
		return fmt.Errorf("message ID is required")
	}
	if message.Type == "" {
		return fmt.Errorf("message type is required")
	}
	if message.Data == nil {
		return fmt.Errorf("message data is required")
	}
	return nil
}
