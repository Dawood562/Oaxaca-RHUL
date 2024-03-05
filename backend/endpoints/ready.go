package endpoints

import "github.com/gofiber/fiber/v2"

func Ready(c *fiber.Ctx) error {
	_, err := GetID(c)
	if err != nil {
		return err
	}
	return nil
}
