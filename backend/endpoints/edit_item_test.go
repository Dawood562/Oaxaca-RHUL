//go:build integration

package endpoints

import (
	"bytes"
	"net/http"
	"teamproject/database"
	"teamproject/database/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestEditItem(t *testing.T) {
	app := fiber.New()
	app.Put("/edit_item", EditItem)

	database.ResetTestMenu()

	testCases := []struct {
		name         string
		json         []byte
		expectedItem models.MenuItem
		code         int
	}{
		{
			name: "WithAllCorrectFields",
			json: []byte(`
				{
					"itemId": 1,
					"itemName": "TESTFOOD5",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: models.MenuItem{
				ID:          1,
				Name:        "TESTFOOD5",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens:   []models.Allergen{},
			},
			code: 200,
		},
		{
			name: "WithChangingAllergens",
			json: []byte(`
				{
					"itemId": 2,
					"itemName": "TESTFOODALLERGENS",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600,
					"allergens": [
						{
							"name": "Gluten"
						},
						{
							"name": "Lactose"
						}
					]
				}
			`),
			expectedItem: models.MenuItem{
				ID:          2,
				Name:        "TESTFOODALLERGENS",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens: []models.Allergen{
					{
						ID:     3,
						ItemID: 2,
						Name:   "Gluten",
					},
					{
						ID:     4,
						ItemID: 2,
						Name:   "Lactose",
					},
				},
			},
			code: 200,
		},
		{
			name: "WithMissingID",
			json: []byte(`
				{
					"itemName": "TESTFOOD6",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: models.MenuItem{
				ID:          1,
				Name:        "TESTFOOD5",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens:   []models.Allergen{},
			},
			code: 422,
		},
		{
			name: "WithInvalidID",
			json: []byte(`
				{
					"itemId": 6,
					"itemName": "TESTFOOD6",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: models.MenuItem{
				ID:          1,
				Name:        "TESTFOOD5",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens:   []models.Allergen{},
			},
			code: 409,
		},
		{
			name: "WithDuplicateName",
			json: []byte(`
				{
					"itemId": 2,
					"itemName": "TESTFOOD5",
					"itemDescription": "New description",
					"price": 6.0,
					"calories": 600
				}
			`),
			expectedItem: models.MenuItem{
				ID:          1,
				Name:        "TESTFOOD5",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens:   []models.Allergen{},
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
			expectedItem: models.MenuItem{
				ID:          1,
				Name:        "TESTFOOD5",
				Description: "New description",
				Price:       6.0,
				Calories:    600,
				Allergens:   []models.Allergen{},
			},
			code: 400,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("PUT", "/edit_item", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			// Send the test request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response status code
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")

			// Check that the item in the database was updated correctly
			menu := database.QueryMenu(&database.MenuFilter{})
			assert.Equal(t, 4, len(menu), "Check that the menu only contains four items")
			assert.Contains(t, menu, test.expectedItem, "Check that the item was correctly updated")
		})
	}
}
