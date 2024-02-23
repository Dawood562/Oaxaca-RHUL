//go:build integration

package endpoints

import (
	"bytes"
	"net/http"
	"teamproject/data"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddWaiters(t *testing.T) {
	app := fiber.New()
	app.Put("/add_waiter", RegisterWaiter)

	testCases := []struct {
		name         string      //test name
		json         []byte      //body to test with
		expectedItem data.Waiter //expected received waiter
		code         int         // expected return
	}{
		{
			name: "TestWaiterWithValidValues",
			json: []byte(`
				{
					"id": 1,
					"waiterUsername": "John",
					"tableNumber": []
				}
			`),
			expectedItem: data.Waiter{ID: 1, Username: "John"},
			code:         200,
		},
		{
			name: "TestWaiterWithInvalidID",
			json: []byte(`
				{
					"id": "1",
					"waiterUsername": "John",
					"tableNumber": []
				}
			`),
			expectedItem: data.Waiter{},
			code:         400,
		},
	}

	for _, test := range testCases {
		data.ClearWaiterList() // Clear waiter list on each test
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("PUT", "/add_waiter", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			// Send the test request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response status code
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")

			// Check waiter list is updated correctly
			waiterData := *data.GetWaiter()
			if test.code <= 299 {
				assert.Equal(t, 1, len(waiterData), "Incorrect number of waiters in waiter list")
				assert.Equal(t, "John", waiterData[0].Username, "Data not inserted correctly")
			}

		})
	}
}
