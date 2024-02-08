//go:build integration

package endpoints

import (
	"testing"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestOpenWebsockets(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		assert.NoError(t, HandleConnection(c))
	}))

	go func() {
		app.Listen(":4444")
	}()
	defer app.Shutdown()

	testCases := []struct {
		name string
		msg  string
		resp string
		clen int
		wlen int
	}{
		{
			name: "CustomerWithCorrectTableNumber",
			msg:  "CUSTOMER:4",
			resp: "WELCOME",
			clen: 1,
			wlen: 0,
		},
		{
			name: "WaiterWithValidName",
			msg:  "WAITER:John",
			resp: "WELCOME",
			clen: 1,
			wlen: 1,
		},
		{
			name: "WaiterWithNoName",
			msg:  "WAITER:",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
		{
			name: "WaiterWithNoNameSeparator",
			msg:  "WAITER",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
		{
			name: "CustomerWithNoTableNumber",
			msg:  "CUSTOMER:",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
		{
			name: "CustomerWithInvalidTableNumber",
			msg:  "CUSTOMER:A",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
		{
			name: "CustomerWithNoTableNumberSeparator",
			msg:  "CUSTOMER",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
		{
			name: "InvalidIdentifier",
			msg:  "TEST",
			resp: "DENIED",
			clen: 1,
			wlen: 1,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a test websocket connection
			retries := 3
			ws, _, err := gwebsocket.DefaultDialer.Dial("ws://127.0.0.1:4444/notifications", nil)
			for retries > 0 && err != nil {
				ws, _, err = gwebsocket.DefaultDialer.Dial("ws://127.0.0.1:4444/notifications", nil)
				retries -= 1
			}
			assert.NoError(t, err, "Test that creating a websocket connection does not create an error")
			defer ws.Close()

			err = ws.WriteMessage(gwebsocket.TextMessage, []byte(test.msg))
			assert.NoError(t, err, "Test sending a message does not create an error")
			_, m, err := ws.ReadMessage()
			assert.NoError(t, err, "Test that receiving a message does not create an error")
			assert.Equal(t, test.resp, string(m), "Test that server replied with expected response")
			assert.Equal(t, test.clen, len(customers), "Test that the number of customers is as expected")
			assert.Equal(t, test.wlen, len(waiters), "Test that the number of waiters is as expected")
		})
	}
}
