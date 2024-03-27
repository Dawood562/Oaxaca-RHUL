package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Pay(c *fiber.Ctx) error {
	// Retrieve required ID arg
	id, err := GetID(c)
	if err != nil {
		return err
	}

	// Attempt to pay for order
	err = database.PayOrder(uint(id))

	if err != nil {
		if err == database.ErrOrderAlreadyPaid {
			return fiber.ErrConflict
		}
		return fiber.ErrNotFound
	}

	BroadcastToWaiters("REFRESH")

	return nil
}
