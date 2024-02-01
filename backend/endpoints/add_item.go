package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

// AddItem is a callback for the API to add an item to the menu
func AddItem(c *fiber.Ctx) error {
	c.Accepts("application/json")
	// Extract the item
	item := database.MenuItem{}
	err := c.BodyParser(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}

	if item.ItemName == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Missing itemName attribute")
	}

	err = database.AddItem(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Item with that name already exists")
	}
	return nil
}
