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

func TestAddItem(t *testing.T) {
	app := fiber.New()
	app.Post("/add_item", AddItem)

	testCases := []struct {
		name              string
		json              []byte
		code              int
		expectedItemNames []string
	}{
		{
			name: "WithCorrectFields",
			json: []byte(`
			{
				"itemName": "TESTFOOD",
				"itemDescription": "A test item",
				"price": 5.00,
				"calories": 500
			}
			`),
			code:              200,
			expectedItemNames: []string{"TESTFOOD"},
		},
		{
			name: "WithDuplicateName",
			json: []byte(`
			{
				"itemName": "TESTFOOD",
				"itemDescription": "A test item",
				"price": 5.00,
				"calories": 500
			}
			`),
			code:              409,
			expectedItemNames: []string{"TESTFOOD"},
		},
		{
			name: "WithMissingItemName",
			json: []byte(`
			{
				"itemDescription": "A test item",
				"price": 5.00,
				"calories": 500
			}
			`),
			code:              422,
			expectedItemNames: []string{"TESTFOOD"},
		},
		{
			name: "WithMissingPrice",
			json: []byte(`
			{
				"itemDescription": "A test item",
				"itemName": "TESTFOOD3",
				"calories": 500
			}
			`),
			code:              200,
			expectedItemNames: []string{"TESTFOOD", "TESTFOOD3"},
		},
		{
			name: "WithMissingCalories",
			json: []byte(`
			{
				"itemDescription": "A test item",
				"itemName": "TESTFOOD4",
				"price": 5.00
			}
			`),
			code:              200,
			expectedItemNames: []string{"TESTFOOD", "TESTFOOD3", "TESTFOOD4"},
		},
		{
			name: "WithMissingDescription",
			json: []byte(`
			{
				"itemName": "TESTFOOD5",
				"price": 5.00,
				"calories": 500
			}
			`),
			code:              200,
			expectedItemNames: []string{"TESTFOOD", "TESTFOOD3", "TESTFOOD4", "TESTFOOD5"},
		},
		{
			name: "WithBadJSON",
			json: []byte(`
			{
				"itemName": "TESTFOOD6",
				"price: 5.00
			}
			`),
			code:              400,
			expectedItemNames: []string{"TESTFOOD", "TESTFOOD3", "TESTFOOD4", "TESTFOOD5"},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("POST", "/add_item", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")
			// Check that the database contains the required entries
			menu := database.QueryMenu(&database.MenuFilter{})
			assert.Equal(t, len(test.expectedItemNames), len(menu), "Check that the database contains the correct number of entries")

			menuNames := make([]string, len(menu))
			for i, item := range menu {
				menuNames[i] = item.ItemName
			}
			for _, expected := range test.expectedItemNames {
				assert.Contains(t, menuNames, expected, "Check that all expected entries are present")
			}
		})
	}

	database.UpdateDB("DELETE FROM menuitem")
}
