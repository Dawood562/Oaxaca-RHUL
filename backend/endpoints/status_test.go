//go:build integration

package endpoints

import (
	"fmt"
	"io"
	"net/http"
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	database.ResetTestOrders()

	app := fiber.New()
	app.Get("/status/:id", Status)

	assertResult(t, "1", app, database.StatusAwaitingConfirmation, fiber.StatusOK)
	assertResult(t, "2", app, database.StatusAwaitingConfirmation, fiber.StatusOK)
	assertResult(t, "3", app, "", fiber.StatusNotFound)
	assertResult(t, "A", app, "", fiber.StatusUnprocessableEntity)
}

// assertResult calls the /status endpoint with the given params and asserts the resulting string response and status code
func assertResult(t *testing.T, id string, app *fiber.App, status string, code int) {
	res, c := callEndpoint(t, id, app)
	if code == fiber.StatusOK {
		assert.Equal(t, status, res, "Test that query returned the expected result")
	}
	assert.Equal(t, code, c, "Test that query returned the expected status code")
}

// callEndpoint is a helper function for calling the /status endpoint and returning the response string and status code
func callEndpoint(t *testing.T, id string, app *fiber.App) (string, int) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/status/%s", id), nil)

	// Send test HTTP request
	res, err := app.Test(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	// Read out response body
	body, _ := io.ReadAll(res.Body)

	return string(body), res.StatusCode
}
