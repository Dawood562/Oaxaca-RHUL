package endpoints

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RemoveWaiterStruct struct {
	ID uint `json:"id"`
}

// RemoveWaiter is an API callback for logging out a waiter
func RemoveWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiterToRemove := RemoveWaiterStruct{}

	bod := string(c.BodyRaw()[:])
	toCompare := `"id":`
	if !strings.Contains(bod, toCompare) {
		return fiber.NewError(fiber.StatusBadRequest, "Request did not include valid/correctly formatted id paramter")
	}
	err := c.BodyParser(&waiterToRemove)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body data")
	}
	err = RemoveWaiterData(WaiterData{ID: waiterToRemove.ID})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Provided id did not match any active waiters")
	}
	return nil
}
