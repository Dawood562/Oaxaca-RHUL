package endpoints

import (
	"strconv"
	"teamproject/database"
	"teamproject/database/models"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {
	tableNumRaw := c.Query("tableNumber")
	status := c.Query("Status")

	tableNum := uint(0)

	if len(tableNumRaw) > 0 {
		temp, err := strconv.ParseUint(tableNumRaw, 10, 32)
		if err != nil {
			c.SendString(err.Error())
		}
		tableNum = uint(temp)
	}

	filter := models.Order{TableNumber: tableNum, Status: status}

	data, err := database.FetchOrders(filter)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.JSON(data)
}
