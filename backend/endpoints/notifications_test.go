//go:build integration

package endpoints

import (
	"strings"
	"testing"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer(t *testing.T) {
	customers = []Customer{}

	testCases := []struct {
		name string
		arg  string
		err  bool
	}{
		{
			name: "WithValidTableNumber",
			arg:  "4",
			err:  false,
		},
		{
			name: "WithDuplicateTableNumber",
			arg:  "4",
			err:  true,
		},
		{
			name: "WithSecondValidTableNumber",
			arg:  "5",
			err:  false,
		},
		{
			name: "WithSecondInvalidTableNumber",
			arg:  "A",
			err:  true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			c, err := createCustomer(test.arg, nil)
			if test.err {
				assert.Error(t, err, "Test that expected error is thrown")
			} else {
				assert.NoError(t, err, "Test that no unexpected errors are thrown")
				customers = append(customers, c.(Customer))
			}
		})
	}
	customers = []Customer{}
}

func TestCreateWaiter(t *testing.T) {

}

func TestRemoveCustomer(t *testing.T) {
	customers = []Customer{{table: 1}, {table: 2}, {table: 3}}
	c := customers[1]
	c.Remove()
	assert.Equal(t, 2, len(customers), "Test that removing a customer really removes the customer")
	c.Remove()
	assert.Equal(t, 2, len(customers), "Test that removing a customer twice does nothing")
	customers[0].Remove()
	assert.Equal(t, 1, len(customers), "Test that removing another customer works correctly")
}

func TestRemoveWaiter(t *testing.T) {

}

func TestOpenAndCloseWebsockets(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		HandleConnection(c)
	}))

	go func() {
		app.Listen(":4444")
	}()
	defer app.Shutdown()

	testCases := []struct {
		name string
		msg  string
		resp string
	}{
		{
			name: "CustomerWithCorrectTableNumber",
			msg:  "CUSTOMER:4",
			resp: "WELCOME",
		},
		{
			name: "WaiterWithValidName",
			msg:  "WAITER:John",
			resp: "WELCOME",
		},
		{
			name: "WaiterWithNoName",
			msg:  "WAITER:",
			resp: "ERROR",
		},
		{
			name: "WaiterWithNoNameSeparator",
			msg:  "WAITER",
			resp: "ERROR",
		},
		{
			name: "CustomerWithNoTableNumber",
			msg:  "CUSTOMER:",
			resp: "ERROR",
		},
		{
			name: "CustomerWithInvalidTableNumber",
			msg:  "CUSTOMER:A",
			resp: "ERROR",
		},
		{
			name: "CustomerWithNoTableNumberSeparator",
			msg:  "CUSTOMER",
			resp: "ERROR",
		},
		{
			name: "InvalidIdentifier",
			msg:  "TEST",
			resp: "ERROR",
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
			assert.True(t, strings.HasPrefix(string(m), test.resp), "Test that the server responded with an expected response")
		})
	}
}
