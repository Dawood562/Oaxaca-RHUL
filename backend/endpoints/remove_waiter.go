package endpoints

import (
	"teamproject/data"

	"github.com/gofiber/fiber/v2"
)

type RemoveWaiterData struct {
	ID uint `json:"id"`
}

func RemoveWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiterToRemove := RemoveWaiterData{}
	err := c.BodyParser(&waiterToRemove)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body data")
	}
	data.RemoveWaiter(data.Waiter{ID: waiterToRemove.ID})
	return nil
}
