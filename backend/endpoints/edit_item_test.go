//go:build integration

package endpoints

import (
	"bytes"
	"net/http"
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestEditItem(t *testing.T) {
	app := fiber.New()
	app.Patch("/edit_item", EditItem)

	// Add an item for editing
	database.UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, itemdescription, price, calories) VALUES (1, 'TESTFOOD', 'Test description', 5.0, 500)")
	database.UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, itemdescription, price, calories) VALUES (2, 'TESTFOOD3', 'Test description', 5.0, 500)")

	testCases := []struct {
		name         string
		json         []byte
		expectedItem database.MenuItem
		code         int
	}{
		{
			name: "WithAllCorrectFields",
			json: []byte(`
				{
					"itemId": 1,
					"itemName": "TESTFOOD2",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: database.MenuItem{
				ID:          1,
				ItemName:    "TESTFOOD2",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
			},
			code: 200,
		},
		{
			name: "WithMissingID",
			json: []byte(`
				{
					"itemName": "TESTFOOD2",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: database.MenuItem{
				ID:          1,
				ItemName:    "TESTFOOD2",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
			},
			code: 422,
		},
		{
			name: "WithInvalidID",
			json: []byte(`
				{
					"itemId": 3,
					"itemName": "TESTFOOD4",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: database.MenuItem{
				ID:          1,
				ItemName:    "TESTFOOD2",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
			},
			code: 409,
		},
		{
			name: "WithDuplicateName",
			json: []byte(`
				{
					"itemId": 2,
					"itemName": "TESTFOOD2",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: database.MenuItem{
				ID:          1,
				ItemName:    "TESTFOOD2",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
			},
			code: 409,
		},
		{
			name: "WithBadJSON",
			json: []byte(`
				{
					itemId: 1,
					itemName: "TESTFOOD3",
					itemDescription: "New description 2",
					price: 6.0,
					calories: 600
				}
			`),
			expectedItem: database.MenuItem{
				ID:          1,
				ItemName:    "TESTFOOD2",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
			},
			code: 400,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("PATCH", "/edit_item", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response status code
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")

			// Check that the item in the database was updated correctly
			menu := database.QueryMenu(&database.MenuFilter{})
			assert.Equal(t, 2, len(menu), "Check that the menu only contains two items")
			assert.Contains(t, menu, test.expectedItem, "Check that the item was correctly updated")
		})
	}

	database.UpdateDB("DELETE FROM menuitem")
}
