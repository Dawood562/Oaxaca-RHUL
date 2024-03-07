package endpoints

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetID extracts the ID parameter from the fiber context. Returns an error if an invalid ID is given.
func GetID(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	if len(idStr) == 0 {
		return 0, fiber.ErrNotFound
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fiber.ErrUnprocessableEntity
	}

	return uint(id), nil
}
