package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func EditItem(c *fiber.Ctx) error {
	c.Accepts("application/json")
	// Extract the item
	item := database.MenuItem{}
	err := c.BodyParser(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}

	// Check that the provided Item ID is valid
	if item.ID == 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Missing required itemId attribute")
	}

	err = database.EditItem(&item)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Could not update item with provided details")
	}

	return nil
}
