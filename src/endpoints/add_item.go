package endpoints

import (
	"teamproject/database"
	"teamproject/database/models"

	"github.com/gofiber/fiber/v2"
)

// AddItem is a callback for the API to add an item to the menu
func AddItem(c *fiber.Ctx) error {
	c.Accepts("application/json")
	// Extract the item
	item := models.MenuItem{}
	err := c.BodyParser(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}
	// Check that an item name was provided
	if item.Name == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Missing itemName attribute")
	}
	err = database.AddItem(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Item with that name already exists")
	}
	return nil
}
