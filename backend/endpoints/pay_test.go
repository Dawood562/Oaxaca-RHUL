//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPay(t *testing.T) {
	app := fiber.New()
	app.Patch("/pay/:id", Pay)

	PatchEndpointTest(t, "pay", app)
}
