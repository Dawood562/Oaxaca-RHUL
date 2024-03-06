//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCancel(t *testing.T) {
	app := fiber.New()
	app.Patch("/cancel/:id", Cancel)

	PatchEndpointTest(t, "cancel", app)
}
