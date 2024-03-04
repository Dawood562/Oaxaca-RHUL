package endpoints

import (
	"strconv"
	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {
	onlyConfirmed := c.Query("confirmed")
	onlyTableNumber := c.Query("tableNumber")

	confirmed := false
	tableNumber := -1

	if len(onlyConfirmed) > 0 {
		temp, err := strconv.ParseBool(onlyConfirmed)
		if err != nil {
			c.SendString(err.Error())
		}
		confirmed = temp
	} else if len(onlyTableNumber) > 0 {
		temp, err := strconv.ParseInt(onlyTableNumber, 10, 32)
		if err != nil {
			c.SendString(err.Error())
		}
		tableNumber = int(temp)
	}

	data, err := database.FetchOrders(confirmed, tableNumber)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(data)
}
