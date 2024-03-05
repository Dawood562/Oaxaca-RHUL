package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Delivered(c *fiber.Ctx) error {
	// Retrieve required ID arg
	idStr := c.Params("id")
	if len(idStr) == 0 {
		return fiber.ErrNotFound
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	// Attempt to mark order as delivered
	err = database.DeliverOrder(uint(id))

	if err != nil {
		if err == database.ErrOrderAlreadyDelivered || err == database.ErrOrderAlreadyCancelled {
			return fiber.ErrConflict
		}
		return fiber.ErrNotFound
	}

	BroadcastToWaiters("REFRESH")
	SendToTable(uint(id), "REFRESH")

	return nil
}
