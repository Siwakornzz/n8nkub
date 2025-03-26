package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Stream bool     `json:"stream"`
	Images []string `json:"images,omitempty"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		log.Printf("Healthcheck requested from %s", c.IP())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/generate", func(c *fiber.Ctx) error {
		log.Printf("Received request from %s: %s", c.IP(), string(c.Body()))

		var userRequest OllamaRequest
		if err := c.BodyParser(&userRequest); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		userRequest.Stream = false

		requestBody, err := json.Marshal(userRequest)
		if err != nil {
			log.Printf("Error encoding request to JSON: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode request"})
		}

		log.Printf("Sending request to Ollama: %s", string(requestBody))

		resp, err := http.Post("http://ollama:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Error calling Ollama API: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to call Ollama"})
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading Ollama response: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response"})
		}

		log.Printf("Received raw response from Ollama: %s", string(body))

		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			log.Printf("Error parsing Ollama response JSON: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse Ollama response"})
		}

		log.Printf("Sending response to client: %s", ollamaResp.Response)

		return c.SendString(ollamaResp.Response)
	})

	// Start Server
	log.Printf("Starting server on :5001")
	log.Fatal(app.Listen(":5001"))
}
