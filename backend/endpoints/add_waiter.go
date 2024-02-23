package endpoints

import (
	"teamproject/data"

	"github.com/gofiber/fiber/v2"
)

// Registers waiter to waiter data structure
func RegisterWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiter := data.Waiter{}
	c.BodyParser(&waiter)
	data.AddWaiter(waiter)
	return nil
}
