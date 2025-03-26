package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// โครงสร้างของ Request ที่จะส่งไปยัง Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // เราจะตั้งค่าเป็น false ตลอด
}

// โครงสร้างของ Response จาก Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
}

func main() {
	app := fiber.New()

	// Healthcheck endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	// API Proxy สำหรับ Ollama
	app.Post("/generate", func(c *fiber.Ctx) error {
		// อ่านข้อมูล JSON จากผู้ใช้
		var userRequest OllamaRequest
		if err := c.BodyParser(&userRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		// ตั้งค่า `stream: false` อัตโนมัติ
		userRequest.Stream = false

		// แปลงเป็น JSON
		requestBody, err := json.Marshal(userRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode request"})
		}

		// ส่ง Request ไปที่ Ollama API
		resp, err := http.Post("http://ollama:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to call Ollama"})
		}
		defer resp.Body.Close()

		// อ่าน Response จาก Ollama
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response"})
		}

		// Parse JSON เพื่อดึงเฉพาะ "response"
		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse Ollama response"})
		}

		// ส่งเฉพาะ "response" กลับไป
		return c.SendString(ollamaResp.Response)
	})

	// Start Server
	log.Fatal(app.Listen(":5001"))
}
