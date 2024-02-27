package endpoints

import (
	"strings"
	"teamproject/data"

	"github.com/gofiber/fiber/v2"
)

type RemoveWaiterData struct {
	ID uint `json:"id"`
}

func RemoveWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiterToRemove := RemoveWaiterData{}

	bod := string(c.BodyRaw()[:])
	toCompare := `"id":`
	if !strings.Contains(bod, toCompare) {
		return fiber.NewError(fiber.StatusBadRequest, "Request did not include valid/correctly formatted id paramter")
	}
	err := c.BodyParser(&waiterToRemove)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body data")
	}
	err = data.RemoveWaiter(data.Waiter{ID: waiterToRemove.ID})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Provided id did not match any active waiters")
	}
	return nil
}
