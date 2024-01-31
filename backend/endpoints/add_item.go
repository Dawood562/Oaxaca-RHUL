package endpoints

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// AddItem is a callback for the API to add an item to the menu
func AddItem(c *fiber.Ctx) error {
	c.Accepts("application/json")
	// Extract the item
	params := c.AllParams()
	fmt.Println(params)
	return nil
}
