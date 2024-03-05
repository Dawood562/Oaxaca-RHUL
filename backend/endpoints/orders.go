package endpoints

import (
	"strconv"
	"teamproject/database"
	"teamproject/database/models"

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
	}

	data, err := database.FetchOrders(confirmed, -1)
	if err != nil {
		return c.SendString(err.Error())
	}

	if len(onlyTableNumber) > 0 {
		temp, err := strconv.ParseInt(onlyTableNumber, 10, 32)
		if err != nil {
			c.SendString(err.Error())
		}
		tableNumber = int(temp)

		filteredData := []models.Order{}
		// Find data for given table number
		for _, dataItem := range data {
			if dataItem.TableNumber == uint(tableNumber) {
				filteredData = append(filteredData, *dataItem)
			}
		}
		return c.JSON(filteredData)
	}

	return c.JSON(data)
}
