package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
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
	// Get the raw body
	body := c.Body()

	// Parse the JSON body manually
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON format",
		})
	}

	// Access the fields from the map
	address, ok := requestBody["address"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Address is required and must be a string",
		})
	}

	message, ok := requestBody["message"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Message is required and must be a string",
		})
	}

	// Create the message for the API key request
	rabbitMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(address),
	}

	if err := globalProducer.PublishMessage(rabbitMessage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	messages, err := globalConsumer.ConsumeMessages()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to consume messages",
		})
	}

	select {
	case msg := <-messages:
		// Process the message directly
		var outerMsg RawMessage
		if err := json.Unmarshal(msg.Body, &outerMsg); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to unmarshal message: %v", err),
			})
		}

		apiKey, err := base64.StdEncoding.DecodeString(outerMsg.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to decode API key: %v", err),
			})
		}

		msg.Ack(false)

		// Make the GPT API call using the API key and original user request
		gptResponse, err := makeGPTCall(string(apiKey), message)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("GPT API call failed: %v", err),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": gptResponse,
		})

	case <-time.After(10 * time.Second):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No API key received within timeout",
		})
	}
}

const apiURL = "https://api.openai.com/v1/chat/completions"

type GPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GPTResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Helper function to make the GPT API call
func makeGPTCall(apiKey string, userRequest string) (string, error) {
	// Prepare the request payload
	requestBody := GPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: userRequest},
		},
	}

	// Convert the request payload to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return "", err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", string(body))
		return "", err
	}

	// Parse the response
	var gptResponse GPTResponse
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return "", err
	}

	// Print the response from GPT
	if len(gptResponse.Choices) > 0 {
		fmt.Println("GPT Response:", gptResponse.Choices[0].Message.Content)
	} else {
		fmt.Println("No response from GPT")
	}

	return gptResponse.Choices[0].Message.Content, nil
}
