//go:build integration

package endpoints

import (
	"testing"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNotifications(t *testing.T) {
	app := fiber.New()

	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		c.WriteMessage(websocket.TextMessage, []byte("Hello Client"))
	}))

	go func() {
		app.Listen(":4444")
	}()
	// Wait for server to start
	time.Sleep(5 * time.Second)

	testCases := []struct {
		name   string
		params map[string]string
	}{
		{
			name: "TemporaryTestCase",
			params: map[string]string{
				"table_num": "5",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a test websocket connection
			ws, _, err := gwebsocket.DefaultDialer.Dial("ws://127.0.0.1:4444/notifications", nil)
			assert.NoError(t, err, "Test that creating a websocket connection does not create an error")
			defer ws.Close()
			// Try reading from the websocket
			_, _, err = ws.ReadMessage()
			assert.NoError(t, err, "Test that receiving a message from the websocket does not create an error")
		})
	}

	app.Shutdown()
}
