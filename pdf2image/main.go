package main

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func ConvertPDF(c *fiber.Ctx) error {
	type Request struct {
		File string `json:"file"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	if req.File == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Field 'file' is required and must contain base64 PDF data"})
	}

	pdfData, err := base64.StdEncoding.DecodeString(req.File)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to decode base64 PDF data"})
	}

	log.Printf("PDF data size: %d bytes", len(pdfData))
	if len(pdfData) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Decoded PDF data is empty"})
	}

	tmpFile, err := ioutil.TempFile("", "pdf2image-*.pdf")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create temporary file"})
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfData); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to write PDF data to temporary file"})
	}
	if err := tmpFile.Close(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to close temporary file"})
	}

	cmd := exec.Command("mutool", "draw", "-F", "png", "-o", "-", "-r", "600", tmpFile.Name())
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		log.Printf("Error running mutool: %v, stderr: %s", err, errBuffer.String())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to convert PDF",
			"details": errBuffer.String(),
		})
	}

	log.Printf("Output size: %d bytes", outBuffer.Len())
	if outBuffer.Len() == 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "No output generated from mutool"})
	}

	// Debug: เซฟภาพ PNG เพื่อตรวจสอบ
	debugFile, err := os.Create("debug_image.png")
	if err != nil {
		log.Printf("Failed to create debug image file: %v", err)
	} else {
		_, err = debugFile.Write(outBuffer.Bytes())
		if err != nil {
			log.Printf("Failed to write debug image: %v", err)
		}
		debugFile.Close()
		log.Println("Debug image saved as debug_image.png")
	}

	img, err := png.Decode(bytes.NewReader(outBuffer.Bytes()))
	if err != nil {
		log.Printf("Error decoding PNG: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode image"})
	}

	var imgBuffer bytes.Buffer
	err = png.Encode(&imgBuffer, img)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode image"})
	}

	base64Str := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())

	return c.JSON(fiber.Map{"message": "Converted successfully", "image": base64Str})
}

func main() {
	app := fiber.New()
	app.Get("/health", HealthCheck)
	app.Post("/convert", ConvertPDF)
	port := "8080"
	log.Printf("🚀 Server is running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
