package endpoints

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required 'file' field")
	}

	// TODO: validation

	// Generate a filename in the format [RANDOM STRING].[EXTENSION]
	s := strings.Split(file.Filename, ".")
	extension := s[len(s)-1]
	body := uuid.New()
	filename := fmt.Sprintf("%s.%s", body.String(), extension)

	err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename))
	if err != nil {
		return err
	}

	return c.SendString(filename)
}
