package endpoints

import "github.com/gofiber/fiber/v2"

func DebugEndpoint(c *fiber.Ctx) error {
	idk := GetQueuedOrders()
	if len(idk) > 0 {
		return c.JSON(idk)
	} else {
		return c.SendString("Queue empty!")
	}

}
