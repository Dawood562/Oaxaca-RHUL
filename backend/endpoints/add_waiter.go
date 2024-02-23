package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

// Registers waiter to waiter data structure
func RegisterWaiter(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
