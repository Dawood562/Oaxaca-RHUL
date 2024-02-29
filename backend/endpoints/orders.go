package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {
	onlyConfirmed := c.Query("confirmed")

	confirmed := false

	if len(onlyConfirmed) > 0 {
		temp, err := strconv.ParseBool(onlyConfirmed)
		if err != nil {
			c.SendString(err.Error())
		}
		confirmed = temp
	}

	data, err := database.FetchOrders(confirmed)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(data)
}
