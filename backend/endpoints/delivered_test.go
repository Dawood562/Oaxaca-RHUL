//go:build integration

package endpoints

import (
	"fmt"
	"net/http"
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDelivered(t *testing.T) {
	app := fiber.New()
	app.Patch("/delivered/:id", Delivered)

	database.ResetTestOrders()

	testCases := []struct {
		name string
		id   string
		code int
	}{
		{
			name: "WithValidID",
			id:   "1",
			code: fiber.StatusOK,
		},
		{
			name: "WithDuplicateID",
			id:   "1",
			code: fiber.StatusConflict,
		},
		{
			name: "WithSecondValidID",
			id:   "2",
			code: fiber.StatusOK,
		},
		{
			name: "WithInvalidID",
			id:   "3",
			code: fiber.StatusNotFound,
		},
		{
			name: "WithNoID",
			id:   "",
			code: fiber.StatusNotFound,
		},
		{
			name: "WithStringID",
			id:   "ABC",
			code: fiber.StatusUnprocessableEntity,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/delivered/%s", test.id), nil)

			// Send test HTTP request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check response code
			assert.Equal(t, test.code, res.StatusCode, "Test that the correct status code was returned")
		})
	}
}
