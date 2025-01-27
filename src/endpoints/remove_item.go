package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

// RemoveItem is an API callback for removing a given item by ID
func RemoveItem(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Query("itemId"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid itemId field")
	}

	err = database.RemoveItem(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}

	return nil
}
