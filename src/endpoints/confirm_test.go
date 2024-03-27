//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestConfirmOrder(t *testing.T) {
	app := fiber.New()
	app.Patch("/confirm/:id", Confirm)

	PatchEndpointTest(t, "confirm", app)
}
