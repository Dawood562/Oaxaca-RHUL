package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

// Status is an API callback for retrieving the status field of an order given by ID
func Status(c *fiber.Ctx) error {
	// Retrieve required ID arg
	id, err := GetID(c)
	if err != nil {
		return err
	}

	status, err := database.GetOrderStatus(uint(id))

	if err != nil {
		return fiber.ErrNotFound
	}

	return c.SendString(status)
}
