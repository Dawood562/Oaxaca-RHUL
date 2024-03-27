// Package endpoints contains a set of handlers that process requests received by each endpoint in the API. The handlers filter data for each endpoint and
// return appropriate HTTP errors if invalid data is given. Any interaction with the database is done in the database package. Handlers return HTTP errors
// if the called database handler returns an error.
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
