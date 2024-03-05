package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Status(c *fiber.Ctx) error {
	// Retrieve required ID arg
	idStr := c.Params("id")
	if len(idStr) == 0 {
		return fiber.ErrNotFound
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	status, err := database.GetOrderStatus(uint(id))

	if err != nil {
		return fiber.ErrNotFound
	}

	return c.SendString(status)
}
