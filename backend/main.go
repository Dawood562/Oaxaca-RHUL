package main

import (
	"teamproject/endpoints"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create new Fiber app
	app := fiber.New()
	// Add logger middleware
	app.Use(logger.New())

	// Register menu route
	app.Get("/menu", endpoints.Menu)

	// Start server
	app.Listen(":4444")
}
