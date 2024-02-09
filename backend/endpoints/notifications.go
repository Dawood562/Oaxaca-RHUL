package endpoints

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

// Users can be removed from the server
type User interface {
	Remove()
}

// Customers are identified by their table numbers
type Customer struct {
	ws    *websocket.Conn
	table uint
}

// Waiters are identified by their names
type Waiter struct {
	ws   *websocket.Conn
	name string
}

var waiters []Waiter
var customers []Customer

// HandleConnection spawns a handler to manage one websocket connection to the server
func HandleConnection(ws *websocket.Conn) {
	u, err := NewConnection(ws)
	if err != nil {
		ws.Close()
		return
	}
	defer ws.Close()

	w, ok := u.(Waiter)
	if ok {
		waiters = append(waiters, w)
	} else {
		c, _ := u.(Customer)
		customers = append(customers, c)
	}

	// Enter a loop to handle further interaction
	for {
		// Receive from websocket
		var m []byte
		if _, m, err = ws.ReadMessage(); err != nil {
			break
		}
		err := HandleMessage(string(m), u)
		if err != nil {
			// TODO: use helpful error message instead
			ws.WriteMessage(websocket.TextMessage, []byte("ERROR"))
		}
	}
	// Cleanup
	u.Remove()
}

// Handle message processed the given message sent by the given user.
// Returns any errors created during handling.
func HandleMessage(m string, u User) error {
	c, ok := u.(Customer)
	if ok {
		// Connection is a customer
		if m == "HELP" {
			r := fmt.Sprintf("HELP:%d", c.table)
			BroadcastToWaiters(r)
			return c.ws.WriteMessage(websocket.TextMessage, []byte("OK"))
		}
	} else {
		_, _ = u.(Waiter)
		// Connection is a waiter
	}

	return errors.New("Invalid command")
}

// BroadcastToWaiters sends m to all connected waiters
func BroadcastToWaiters(m string) {
	for _, w := range waiters {
		if w.ws != nil {
			w.ws.WriteMessage(websocket.TextMessage, []byte(m))
		}
	}
}

// NewConnection registers a new connection to the system given its websocket connection.
// Returns the User that has been registered.
// Returns an error if one arises during communication.
func NewConnection(ws *websocket.Conn) (User, error) {
	_, m, err := ws.ReadMessage()
	if err != nil {
		return Customer{}, err
	}
	c, err := registerConnection(string(m), ws)
	if err != nil {
		return Customer{}, ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ERROR: %s", err.Error())))
	}
	return c, ws.WriteMessage(websocket.TextMessage, []byte("WELCOME"))
}

// registerConnection reads the first message from the websocket and registers the connection. ws is stored for later use.
// A useful error message is returned if the information given is incomplete.
// Returns true if the connection is a customer.
func registerConnection(m string, ws *websocket.Conn) (User, error) {
	segments := strings.Split(string(m), ":")
	if len(segments) != 2 {
		return nil, errors.New("Invalid payload format")
	}

	var ret User
	var err error
	switch segments[0] {
	case "CUSTOMER":
		ret, err = createCustomer(segments[1], ws)
	case "WAITER":
		ret, err = createWaiter(segments[1], ws)
	default:
		err = errors.New("Invalid identifier given")
	}
	return ret, err
}

// createCustomer creates a customer with the given websocket with the given table number.
// An error is returned if a customer with the same table number already exists with an open connection.
func createCustomer(arg string, ws *websocket.Conn) (User, error) {
	n, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		return Customer{}, errors.New("Table number must be a number")
	}
	for _, c := range customers {
		if int(c.table) == int(n) {
			return Customer{}, fmt.Errorf("Table number %d is already connected", n)
		}
	}
	return Customer{ws: ws, table: uint(n)}, nil
}

// Remove removes this customer from the server register
func (c Customer) Remove() {
	for i, cu := range customers {
		if cu.table == c.table {
			// Remove the customer
			customers[i] = customers[len(customers)-1]
			customers = customers[:len(customers)-1]
		}
	}
}

// createWaiter creates a waiter connection with the given websocket with the given name
func createWaiter(arg string, ws *websocket.Conn) (User, error) {
	if len(arg) == 0 {
		return Waiter{}, errors.New("Waiter name not provided")
	}
	for _, c := range waiters {
		if c.name == arg {
			return Waiter{}, errors.New("Waiter with that name is already connected")
		}
	}
	return Waiter{ws: ws, name: arg}, nil
}

// Remove removes a waiter from the server register
func (w Waiter) Remove() {
	for i, wa := range waiters {
		if wa.name == w.name {
			// Remove the waiter
			waiters[i] = waiters[len(waiters)-1]
			waiters = waiters[:len(waiters)-1]
		}
	}
}
