package endpoints

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

type WaiterAPI struct {
	Username string `json:"username"`
}

// Registers waiter to waiter data structure
func RegisterWaiter(c *fiber.Ctx) error {
	c.Accepts("application/json")
	waiter := WaiterAPI{}
	err := c.BodyParser(&waiter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}

	// Generate random id
	randWaiterId := uint(rand.Intn(100000000))
	waiterData := WaiterData{ID: randWaiterId, Username: waiter.Username}
	err = AddWaiterData(waiterData)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Cannot add waiter with same id as existing waiter")
	}
	return nil
}
