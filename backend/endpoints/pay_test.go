//go:build integration

package endpoints

import (
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPay(t *testing.T) {
	app := fiber.New()
	app.Put("/pay", Pay)
	database.ResetTestMenu()
}
