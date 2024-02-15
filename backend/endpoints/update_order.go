package endpoints

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RemoveOrder(c *fiber.Ctx) error {
	_, err := strconv.ParseInt(c.Query("itemId"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid itemId field")
	}

	return c.SendString("OK")
}
