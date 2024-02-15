package endpoints

import (
	"fmt"
	"teamproject/database/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {

	exampleItems := []models.MenuItem{models.MenuItem{ID: 1, Name: "Tequila", Description: "the good stuff", Price: 22.99, Calories: 753},
		models.MenuItem{ID: 2, Name: "Vodka", Description: "the wake up in a bush stuff", Price: 16.99, Calories: 713}}

	exampleOrderItems := []models.OrderItem{models.OrderItem{OrderID: 1, Item: exampleItems[0], Notes: "Good order"},
		models.OrderItem{OrderID: 2, Item: exampleItems[1], Notes: "Questionable order"}}

	exampleOrder := models.Order{
		ID:          1,
		Time:        time.Now(),
		TableNumber: 13,
		Bill:        6.99,
		Status:      "Ready",
		Items:       exampleOrderItems,
	}

	fmt.Println(exampleOrder)

	return c.JSON(exampleOrder)
}
