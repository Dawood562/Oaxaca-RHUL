package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Cancel(c *fiber.Ctx) error {
	// Retrieve required ID arg
	id, err := GetID(c)
	if err != nil {
		return err
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
	BroadcastToWaiters("REFRESH")
	BroadcastToCustomers("REFRESH")

	return nil
}
