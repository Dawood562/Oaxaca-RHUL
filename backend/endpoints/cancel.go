package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Cancel(c *fiber.Ctx) error {
	// Retrieve required ID arg
	idStr := c.Params("id")
	if len(idStr) == 0 {
		return fiber.ErrNotFound
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	// Attempt to pay for order
	err = database.CancelOrder(uint(id))

	if err != nil {
		if err == database.ErrOrderAlreadyCancelled {
			return fiber.ErrConflict
		}
		return fiber.ErrNotFound
	}

	BroadcastToKitchen("CANCEL")

	return nil
}
