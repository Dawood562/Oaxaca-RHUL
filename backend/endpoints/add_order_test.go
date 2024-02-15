//go:build integration

package endpoints

import (
	"bytes"
	"io"
	"net/http"
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	database.ClearOrders()
	database.ResetTestMenu()

	app := fiber.New()
	app.Post("/add_order", AddOrder)

	testCases := []struct {
		name  string
		json  []byte
		code  int
		tnum  int   // The table number to check order items for s
		count int   // Expected count of orders in database
		items []int // Expected menu item ID list
	}{
		{
			name: "WithCorrectFields",
			json: []byte(`
			{
				"tableNumber": 4,
				"bill": 50.0,
				"items": [
					{
						"item": 1,
						"notes": "Gluten free"
					},
					{
						"item": 2,
						"notes": ""
					}
				]
			}
			`),
			code:  200,
			tnum:  4,
			count: 1,
			items: []int{1, 2},
		},
		{
			name: "DuplicateOrder",
			json: []byte(`
			{
				"tableNumber": 4,
				"bill": 50.0,
				"items": [
					{
						"item": 1,
						"notes": "Gluten free"
					},
					{
						"item": 2,
						"notes": ""
					}
				]
			}
			`),
			code:  409,
			tnum:  4,
			count: 1,
			items: []int{},
		},
		{
			name: "SecondValidOrder",
			json: []byte(`
			{
				"tableNumber": 5,
				"bill": 50.0,
				"items": [
					{
						"item": 3,
						"notes": "Gluten free"
					},
					{
						"item": 4,
						"notes": ""
					}
				]
			}
			`),
			code:  200,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
		{
			name: "NoItems",
			json: []byte(`
			{
				"tableNumber": 5,
				"bill": 50.0
			}
			`),
			code:  422,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
		{
			name: "InvalidItemID",
			json: []byte(`
			{
				"tableNumber": 6,
				"bill": 50.0,
				"items": [
					{
						"item": 10,
						"notes": "Gluten free"
					},
					{
						"item": 4,
						"notes": ""
					}
				]
			}
			`),
			code:  409,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
		{
			name: "NoItemID",
			json: []byte(`
			{
				"tableNumber": 6,
				"bill": 50.0,
				"items": [
					{
						"notes": "Gluten free"
					},
					{
						"item": 4,
						"notes": ""
					}
				]
			}
			`),
			code:  409,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
		{
			name: "NoTableNumber",
			json: []byte(`
			{
				"bill": 50.0,
				"items": [
					{
						"item": 3,
						"notes": "Gluten free"
					},
					{
						"item": 4,
						"notes": ""
					}
				]
			}
			`),
			code:  422,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
		{
			name: "InvalidJSON",
			json: []byte(`
			{
				"tableNumber": 6,
				"bill": 50.0
				"items": [
					{
						"item": 3,
						"notes": "Gluten free"
					},
					{
						"item": 4,
						"notes": ""
					}
				]
			}
			`),
			code:  400,
			tnum:  5,
			count: 2,
			items: []int{3, 4},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("POST", "/add_order", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			// Send test HTTP request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.code, res.StatusCode, string(body))
			// Check the number of orders in the database
			assert.Equal(t, test.count, database.CountOrders(), "Check that the database contains the correct number of orders")
			// Check item IDs
			if test.code == 200 {
				ids := database.FetchOrderItemIDs(test.tnum)
				// Check that all IDs are present
				for _, id := range test.items {
					assert.Contains(t, ids, id, "Test that the database contains the expected items")
				}
			}
		})
	}
}
