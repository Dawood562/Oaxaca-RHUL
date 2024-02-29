package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Confirm(c *fiber.Ctx) error {
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
	err = database.ConfirmOrder(uint(id))

	if err != nil {
		if err == database.ErrOrderAlreadyConfirmed {
			return fiber.ErrConflict
		}
		return fiber.ErrNotFound
	}

	BroadcastToKitchen("CONFIRM")
	BroadcastToWaiters("REFRESH")

	return nil
}
