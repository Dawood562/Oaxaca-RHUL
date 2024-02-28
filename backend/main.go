package main

import (
	"teamproject/endpoints"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create new Fiber app
	app := fiber.New()
	// Add logger and cors middlewares
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./swagger.json",
		Path:     "docs",
		Title:    "Team Project API",
	}))

	// Register routes
	app.Get("/menu", endpoints.Menu)
	app.Post("/add_item", endpoints.AddItem)
	app.Delete("/remove_item", endpoints.RemoveItem)
	app.Put("/edit_item", endpoints.EditItem)

	app.Post("/add_order", endpoints.AddOrder)
	app.Delete("/remove_order", endpoints.RemoveOrder)
	app.Get("/orders", endpoints.GetOrders)
	app.Patch("/pay/:id", endpoints.Pay)
	app.Patch("/confirm/:id", endpoints.Confirm)
	app.Patch("/cancel/:id", endpoints.Cancel)

	// Register websocket endpoint
	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		endpoints.HandleConnection(c)
	}))

	// Start server
	app.Listen(":4444")
}
