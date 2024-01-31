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
		name string
		json []byte
		code int
	}{
		{
			name: "WithCorrectFields",
			json: []byte(`
			{
				"itemName": "TESTFOOD",
				"price": 5.00,
				"calories": 500
			}
			`),
			code: 200,
		},
		{
			name: "WithDuplicateName",
			json: []byte(`
			{
				"itemName": "TESTFOOD",
				"price": 5.00,
				"calories": 500
			}
			`),
			code: 409,
		},
		{
			name: "WithMissingItemName",
			json: []byte(`
			{
				"price": 5.00,
				"calories": 500
			}
			`),
			code: 422,
		},
		{
			name: "WithMissingPrice",
			json: []byte(`
			{
				"itemName": "TESTFOOD3",
				"calories": 500
			}
			`),
			code: 200,
		},
		{
			name: "WithMissingCalories",
			json: []byte(`
			{
				"itemName": "TESTFOOD4",
				"price": 5.00
			}
			`),
			code: 200,
		},
		{
			name: "WithBadJSON",
			json: []byte(`
			{
				"itemName": "TESTFOOD5",
				"price: 5.00
			}
			`),
			code: 400,
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
		})
	}

	database.UpdateDB("DELETE FROM menuitem")
}
