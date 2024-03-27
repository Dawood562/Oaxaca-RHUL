//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestReady(t *testing.T) {
	app := fiber.New()
	app.Patch("/ready/:id", Ready)

	PatchEndpointTest(t, "ready", app)
}
