package main

import (
	"teamproject/endpoints"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create new Fiber app
	app := fiber.New()
	// Add logger middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./swagger.json",
		Path:     "docs",
		Title:    "Team Project API",
	}))

	// Register menu route
	app.Get("/menu", endpoints.Menu)

	// Start server
	app.Listen(":4444")
}
