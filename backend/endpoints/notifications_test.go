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
	customers.users = make([]User, 0)

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
				customers.Append(c)
			}
		})
	}
}

func TestCreateWaiter(t *testing.T) {
	waiters.users = make([]User, 0)

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
				waiters.Append(w)
			}
		})
	}
}

func TestRemoveCustomer(t *testing.T) {
	customers.users = []User{Customer{table: 1}, Customer{table: 2}, Customer{table: 3}}
	c := customers.users[1]
	c.Remove()
	assert.Equal(t, 2, len(customers.users), "Test that removing a customer really removes the customer")
	c.Remove()
	assert.Equal(t, 2, len(customers.users), "Test that removing a customer twice does nothing")
	customers.users[0].Remove()
	assert.Equal(t, 1, len(customers.users), "Test that removing another customer works correctly")
}

func TestRemoveWaiter(t *testing.T) {
	waiters.users = []User{Waiter{name: "Jacob"}, Waiter{name: "Josh"}, Waiter{name: "Dawood"}}
	w := waiters.users[1]
	w.Remove()
	assert.Equal(t, 2, len(waiters.users), "Test that removing a waiter really removes the customer")
	w.Remove()
	assert.Equal(t, 2, len(waiters.users), "Test that removing a waiter twice does nothing")
	waiters.users[0].Remove()
	assert.Equal(t, 1, len(waiters.users), "Test that removing another waiter works correctly")
}

func TestRemoveKitchen(t *testing.T) {
	kitchens.users = []User{Kitchen{}, Kitchen{}, Kitchen{}}
	k := kitchens.users[1]
	k.Remove()
	assert.Equal(t, 2, len(kitchens.users), "Test that removing a kitchen really removes the customer")
	kitchens.users[0].Remove()
	assert.Equal(t, 1, len(kitchens.users), "Test that removing another kitchen works correctly")
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
			name: "Kitchen",
			msg:  "KITCHEN",
			resp: "WELCOME",
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
	customers.users = []User{}
	waiters.users = []User{}

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

	// Initialise customer connections
	csockets := []*gwebsocket.Conn{createTestWebsocket(t), createTestWebsocket(t)}
	for i, c := range csockets {
		defer c.Close()
		err := c.WriteMessage(gwebsocket.TextMessage, []byte(fmt.Sprintf("CUSTOMER:%d", i+1)))
		assert.NoError(t, err, "Test that sending initial message creates no errors")
		_, m, err := c.ReadMessage()
		assert.Equal(t, "WELCOME", string(m), "Test that server accepted connection and authentication")
	}
	// Initialize waiter connections
	wsockets := []*gwebsocket.Conn{createTestWebsocket(t), createTestWebsocket(t)}
	for i, w := range wsockets {
		defer w.Close()
		err := w.WriteMessage(gwebsocket.TextMessage, []byte(fmt.Sprintf("WAITER:%d", i+1)))
		assert.NoError(t, err, "Test that sending initial message creates no errors")
		_, m, err := w.ReadMessage()
		assert.Equal(t, "WELCOME", string(m), "Test that server accepted connection and authentication")
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Use specified client to send test message
			err := csockets[test.cid].WriteMessage(gwebsocket.TextMessage, []byte(test.cmsg))
			assert.NoError(t, err, "Test that sending help message creates no errors")
			// Read response from server
			_, m, err := csockets[test.cid].ReadMessage()
			assert.NoError(t, err, "Test that receiving response from server creates no errors")
			assert.Equal(t, test.resp, string(m), "Test that the server gave the expected response")
			// Check what the waiters received
			for i, w := range wsockets {
				if len(test.wrecv) > 0 {
					_, m, err = w.ReadMessage()
					assert.NoError(t, err, "Test that reading from waiter socket creates no errors")
					assert.Equal(t, test.wrecv, string(m), fmt.Sprintf("Test that waiter %d received the expected message", i))
				}
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
