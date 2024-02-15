package endpoints

import (
	"teamproject/database/models"

	"github.com/gofiber/fiber/v2"
)

// Boiler plate code for AddOrder while db content isnt
func AddOrder(c *fiber.Ctx) error {
	c.Accepts("application/json")
	item := models.Order{}
	err := c.BodyParser(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}
	return c.SendString("OK")
}
