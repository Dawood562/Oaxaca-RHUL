package endpoints

import (
	"fmt"
	"teamproject/database"
	"teamproject/database/models"

	"github.com/gofiber/fiber/v2"
)

type APIOrder struct {
	TableNumber uint    `json:"tableNumber"`
	Bill        float64 `json:"bill"`
	Items       []struct {
		Item  int    `json:"item"`
		Notes string `json:"notes"`
	} `json:"items"`
}

// Boiler plate code for AddOrder while db content isnt
func AddOrder(c *fiber.Ctx) error {
	c.Accepts("application/json")
	order := APIOrder{}
	err := c.BodyParser(&order)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not process provided JSON")
	}
	// Check that required fields were provided
	if order.TableNumber <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Missing required tableNumber field")
	}
	if len(order.Items) == 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Order has no items")
	}
	// Turn provided data into an order
	o, err := createOrder(&order)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}
	err = database.AddOrder(&o)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}

	// Broadcast notification to waiters
	BroadcastToWaiters("NEW")

	// Send the client their order ID
	c.SendString(fmt.Sprintf("%d", o.ID))

	return nil
}

func createOrder(o *APIOrder) (models.Order, error) {
	apiItems := o.Items
	items := make([]models.OrderItem, len(apiItems))
	for i, item := range apiItems {
		// Retrieve item from database
		ditem, err := database.FetchItem(item.Item)
		if err != nil {
			return models.Order{}, err
		}
		// Convert into an Order Item
		oitem := models.OrderItem{
			Item:  ditem,
			Notes: item.Notes,
		}
		items[i] = oitem
	}

	// Create the Order
	order := models.Order{
		TableNumber: o.TableNumber,
		Bill:        o.Bill,
		Status:      "Awaiting Confirmation",
		Items:       items,
	}

	return order, nil
}
