package endpoints

import (
	"encoding/json"

	"teamproject/database"

	"github.com/gofiber/fiber/v2"
)

// Menu retrieves the menu from the database and represents it as a JSON string.
func Menu(c *fiber.Ctx) error {
	r, err := FetchMenu(nil)
	if err != nil {
		return err
	}
	return c.SendString(r)
}

// FetchMenu returns the current restaurant menu as a JSON encoded string
func FetchMenu(_ *database.MenuFilter) (string, error) {
	r, err := json.Marshal(database.QueryMenu())
	return string(r), err
}
