//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestDelivered(t *testing.T) {
	app := fiber.New()
	app.Patch("/delivered/:id", Delivered)

	PatchEndpointTest(t, "delivered", app)
}
