package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Delivered(c *fiber.Ctx) error {
	// Retrieve required ID arg
	id, err := GetID(c)
	if err != nil {
		return err
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
