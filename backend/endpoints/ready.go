package endpoints

import (
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func Ready(c *fiber.Ctx) error {
	id, err := GetID(c)
	if err != nil {
		return err
	}

	err = database.ReadyOrder(id)

	if err != nil {
		if err == database.ErrOrderNotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrConflict
	}

	BroadcastToWaiters("REFRESH")

	return nil
}
