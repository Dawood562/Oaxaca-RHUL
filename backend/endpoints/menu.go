package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

// Menu retrieves the menu from the database and represents it as a JSON string.
func Menu(c *fiber.Ctx) error {
	// Placeholder dummy return to test API
	return c.SendString("This is a menu")
}
