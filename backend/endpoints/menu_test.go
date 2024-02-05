//go:build integration

package endpoints

import (
	"encoding/json"
	"net/http"
	"teamproject/database"
	"teamproject/database/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	// Setup test data
	database.ResetTestMenu()

	// Setup server for testing
	app := fiber.New()
	app.Get("/menu", Menu)

	testCases := []struct {
		name  string
		args  map[string]string
		items []string
	}{
		{
			name:  "NoFilters",
			args:  map[string]string{},
			items: []string{"TESTFOOD", "TESTFOOD2", "TESTFOOD3", "TESTFOOD4"},
		},
		{
			name:  "WithSearchTerm",
			args:  map[string]string{"searchTerm": "FOOD2"},
			items: []string{"TESTFOOD2"},
		},
		{
			name:  "WithMaxPrice",
			args:  map[string]string{"maxPrice": "5.50"},
			items: []string{"TESTFOOD"},
		},
		{
			name:  "WithMaxCalories",
			args:  map[string]string{"maxCalories": "500"},
			items: []string{"TESTFOOD", "TESTFOOD2"},
		},
		{
			name:  "WithMultipleFilters",
			args:  map[string]string{"searchTerm": "TESTFOOD", "maxPrice": "6.50", "maxCalories": "450"},
			items: []string{"TESTFOOD"},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("GET", "/menu", nil)

			q := req.URL.Query()
			// Add args
			for k, v := range test.args {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			// Run the query
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response
			assert.Equal(t, 200, res.StatusCode, "Check that request returned a 200 OK status code")
			// Check that the returned response is JSON
			var data []models.MenuItem
			err = json.NewDecoder(res.Body).Decode(&data)
			assert.NoError(t, err, "Check that the request returned valid JSON")

			// Check that expected items matches the items returned
			assert.Equal(t, len(test.items), len(data), "Test that the right number of items were returned")
			names := make([]string, len(data))
			for i, item := range data {
				names[i] = item.Name
			}
			for _, item := range test.items {
				assert.Contains(t, names, item, "Check that all required items were returned")
			}
		})
	}
}
