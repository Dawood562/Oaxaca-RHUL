package endpoints

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Upload is an API callback for uploading an image file for use in the menu
func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required 'file' field")
	}

	// Check file type
	data, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Cannot process given file")
	}
	content, err := io.ReadAll(data)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Could not read image")
	}
	t := http.DetectContentType(content)
	if t != "image/png" && t != "image/jpeg" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file type")
	}

	// Generate a filename in the format [RANDOM STRING].[EXTENSION]
	s := strings.Split(file.Filename, ".")
	extension := strings.ToLower(s[len(s)-1])
	// Check file extension
	if extension != "png" && extension != "jpg" && extension != "jpeg" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file extension")
	}
	body := uuid.New()
	filename := fmt.Sprintf("%s.%s", body.String(), extension)

	err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename))
	if err != nil {
		return err
	}

	return c.SendString(filename)
}
