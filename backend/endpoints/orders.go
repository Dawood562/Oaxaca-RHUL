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
	providedId := c.Query("waiterId")

	confirmed := false
	tableNumber := -1

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

	if len(providedId) > 0 {
		waiterID, err := strconv.ParseInt(providedId, 10, 32)
		if err != nil {
			return c.SendString(err.Error())
		}
		// Find waiters with same id
		waiterList := GetWaiter(WaiterData{ID: uint(waiterID)})
		if len(*waiterList) <= 0 {
			return c.SendString("Error finding waiter from provided id")
		}

		waiter := (*waiterList)[0]
		filteredData := []models.Order{}
		// Find data for given table number
		for _, dataItem := range data {
			for _, waiterTableNum := range waiter.TableNumber {
				// If table allocated to waiter found add it to return list
				if dataItem.TableNumber == waiterTableNum {
					filteredData = append(filteredData, *dataItem)
				}
			}
		}
		return c.JSON(filteredData)
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
