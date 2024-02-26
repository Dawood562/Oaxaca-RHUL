package endpoints

import (
	"teamproject/data"

	"github.com/gofiber/fiber/v2"
)

// Registers waiter to waiter data structure
func RegisterWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiter := data.Waiter{}
	err := c.BodyParser(&waiter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}
	data.AddWaiter(waiter)
	return nil
}
