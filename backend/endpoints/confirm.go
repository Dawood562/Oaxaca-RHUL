package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Confirm(c *fiber.Ctx) error {
	// Retrieve required ID arg
	id, err := GetID(c)
	if err != nil {
		return err
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
