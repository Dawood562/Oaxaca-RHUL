package endpoints

import (
	"encoding/json"
	"strconv"

	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

// Menu retrieves the menu from the database and represents it as a JSON string.
func Menu(c *fiber.Ctx) error {
	// Check provided params
	maxPrice, err := strconv.ParseFloat(c.Query("maxPrice"), 32)
	if err != nil {
		maxPrice = 0
	}
	maxCalories, err := strconv.ParseInt(c.Query("maxCalories"), 10, 32)
	if err != nil {
		maxCalories = 0
	}
	filter := &database.MenuFilter{
		SearchTerm:  c.Query("searchTerm"),
		MaxPrice:    float32(maxPrice),
		MaxCalories: int(maxCalories),
	}
	// Retrieve the menu
	r, err := FetchMenu(filter)
	if err != nil {
		return err
	}
	return c.SendString(r)
}

// FetchMenu returns the current restaurant menu as a JSON encoded string.
// filter is applied to the database query to filter results
func FetchMenu(filter *database.MenuFilter) (string, error) {
	r, err := json.Marshal(database.QueryMenu(filter))
	return string(r), err
}
