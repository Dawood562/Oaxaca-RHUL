package endpoints

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRemoveWaiters(t *testing.T) {
	app := fiber.New()
	app.Post("/remove_waiter", RemoveWaiter)

	testCases := []struct {
		name string //test name
		json []byte //body to test with
		code int    // expected return
	}{
		{
			name: "TestRemoveWaiterWithCorrectParams",
			json: []byte(`
			{
				"id": 1
			}
			`),
			code: 200,
		},
		{
			name: "TestRemoveWaiterWithIncorrectData",
			json: []byte(`
			{
				"id": "1"
			}
			`),
			code: 400,
		},
		{
			name: "TestRemoveWaiterWithNoID",
			json: []byte(`
			{
			}
			`),
			code: 400,
		},
		{
			name: "TestRemoveWaiterWithNoInvalidID",
			json: []byte(`
			{
				"id":69
			}
			`),
			code: 400,
		},
	}

	for _, test := range testCases {
		// Clear waiter data and then add test data
		ClearWaiterList()
		AddWaiterData(WaiterData{ID: 1, Username: "John"})
		AddWaiterData(WaiterData{ID: 2, Username: "Bob"})

		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("POST", "/remove_waiter", bytes.NewBuffer(test.json))
			req.Header.Set("Content-Type", "application/json")

			// Send the test request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response status code
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")

			// Check waiter list is updated correctly
			waiterData := *GetWaiter()
			if test.code <= 299 { // If the operation was successful
				assert.Equal(t, 1, len(waiterData), "Incorrect number of waiters in waiter list")
				assert.Equal(t, "Bob", waiterData[0].Username, "Data not inserted correctly")
			}

		})
	}
}
