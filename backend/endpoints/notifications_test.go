//go:build integration

package endpoints

import (
	"fmt"
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
	waiters = []Waiter{}

	testCases := []struct {
		name string
		arg  string
		err  bool
	}{
		{
			name: "WithValidName",
			arg:  "Jacob",
			err:  false,
		},
		{
			name: "WithDuplicateName",
			arg:  "Jacob",
			err:  true,
		},
		{
			name: "WithNoName",
			arg:  "",
			err:  true,
		},
		{
			name: "WithSecondValidName",
			arg:  "Josh",
			err:  false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			w, err := createWaiter(test.arg, nil)
			if test.err {
				assert.Error(t, err, "Test that expected error is thrown")
			} else {
				assert.NoError(t, err, "Test that no unexpected errors are thrown")
				waiters = append(waiters, w.(Waiter))
			}
		})
	}
	waiters = []Waiter{}
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
	waiters = []Waiter{{name: "Jacob"}, {name: "Josh"}, {name: "Dawood"}}
	w := waiters[1]
	w.Remove()
	assert.Equal(t, 2, len(waiters), "Test that removing a waiter really removes the customer")
	w.Remove()
	assert.Equal(t, 2, len(waiters), "Test that removing a waiter twice does nothing")
	waiters[0].Remove()
	assert.Equal(t, 1, len(waiters), "Test that removing another waiter works correctly")
}

func TestOpenAndCloseWebsockets(t *testing.T) {
	app := createTestServer()
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
			ws := createTestWebsocket(t)
			defer ws.Close()

			err := ws.WriteMessage(gwebsocket.TextMessage, []byte(test.msg))
			assert.NoError(t, err, "Test sending a message does not create an error")
			_, m, err := ws.ReadMessage()
			assert.NoError(t, err, "Test that receiving a message does not create an error")
			assert.True(t, strings.HasPrefix(string(m), test.resp), "Test that the server responded with an expected response")
		})
	}
}

func TestCustomerCallWaiter(t *testing.T) {
	app := createTestServer()
	defer app.Shutdown()

	testCases := []struct {
		name  string
		cid   int    // Index of customer to send help call
		cmsg  string // Message sent from customer to server
		resp  string // Expected response from server
		wrecv string // Expected message received by waiter socket
	}{
		{
			name:  "ValidHelpCall",
			cid:   0,
			cmsg:  "HELP",
			resp:  "OK",
			wrecv: "HELP:1",
		},
		{
			name:  "ValidHelpCall",
			cid:   1,
			cmsg:  "HELP",
			resp:  "OK",
			wrecv: "HELP:2",
		},
		{
			name:  "InvalidMessage",
			cid:   0,
			cmsg:  "TEST",
			resp:  "ERROR",
			wrecv: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Initialise customer connections
			customers := []*gwebsocket.Conn{createTestWebsocket(t), createTestWebsocket(t)}
			for i, c := range customers {
				defer c.Close()
				err := c.WriteMessage(gwebsocket.TextMessage, []byte(fmt.Sprintf("CUSTOMER:%d", i+1)))
				assert.NoError(t, err, "Test that sending initial message creates no errors")
				_, m, err := c.ReadMessage()
				assert.Equal(t, "WELCOME", string(m), "Test that server accepted connection and authentication")
			}
		})
	}
}

// createTestWebsocket creates a test websocket connection and returns it
func createTestWebsocket(t *testing.T) *gwebsocket.Conn {
	retries := 3
	ws, _, err := gwebsocket.DefaultDialer.Dial("ws://127.0.0.1:4444/notifications", nil)
	for retries > 0 && err != nil {
		ws, _, err = gwebsocket.DefaultDialer.Dial("ws://127.0.0.1:4444/notifications", nil)
		retries -= 1
	}
	assert.NoError(t, err, "Test that creating a websocket connection does not create an error")
	return ws
}

// createTestServer creates and returns a Fiber server for testing
func createTestServer() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/notifications", websocket.New(func(c *websocket.Conn) {
		HandleConnection(c)
	}))

	go func() {
		app.Listen(":4444")
	}()
	return app
}
