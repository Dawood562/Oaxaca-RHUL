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
	// Menu interaction
	app.Get("/menu", endpoints.Menu)
	app.Post("/add_item", endpoints.AddItem)
	app.Delete("/remove_item", endpoints.RemoveItem)
	app.Put("/edit_item", endpoints.EditItem)
	// Order interaction
	app.Post("/add_order", endpoints.AddOrder)
	app.Get("/orders", endpoints.GetOrders)
	app.Patch("/pay/:id", endpoints.Pay)
	app.Put("/add_waiter", endpoints.RegisterWaiter)
	app.Post("/remove_waiter", endpoints.RemoveWaiter)
	// Order tracking
	app.Patch("/confirm/:id", endpoints.Confirm)
	app.Patch("/cancel/:id", endpoints.Cancel)
	app.Patch("/pay/:id", endpoints.Pay)
	app.Patch("/delivered/:id", endpoints.Delivered)
	app.Get("/status/:id", endpoints.Status)
	app.Patch("/ready/:id", endpoints.Ready)
	// Image uploading and serving
	app.Post("/upload", endpoints.Upload)
	app.Static("/image", "./uploads")

	// Live websocket notifications
	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		endpoints.HandleConnection(c)
	}))

	// Start server
	app.Listen(":4444")
}
