package endpoints

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

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

// Kitchen staff are anonymous
type Kitchen struct {
	ws *websocket.Conn
}

// Registry is a thread-safe slice type
type Registry struct {
	sync.Mutex
	users []User
}

var waiters Registry
var customers Registry
var kitchens Registry

// HandleConnection spawns a handler to manage one websocket connection to the server
func HandleConnection(ws *websocket.Conn) {
	u, err := NewConnection(ws)
	if err != nil {
		ws.Close()
		return
	}
	defer ws.Close()

	// Sort the created connection into the correct registry
	w, ok := u.(Waiter)
	if ok {
		waiters.Append(w)
	} else {
		c, ok := u.(Customer)
		if ok {
			customers.Append(c)
		} else {
			kitchens.Append(u)
		}
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
		} else if m == "NEW" {
			BroadcastToWaiters("NEW")
			return c.ws.WriteMessage(websocket.TextMessage, []byte("OK"))
		}
	} else {
		_, ok = u.(Waiter)
		if ok {
			// Connection is a waiter
		} else {
			// Connection is kitchen staff
			if m == "SERVICE" {
				BroadcastToWaiters("SERVICE")
				return u.(Kitchen).ws.WriteMessage(websocket.TextMessage, []byte("OK"))
			}
		}
	}

	return errors.New("Invalid command")
}

// BroadcastToWaiters sends m to all connected waiters
func BroadcastToWaiters(m string) {
	waiters.Lock()
	defer waiters.Unlock()
	for _, w := range waiters.users {
		w := w.(Waiter)
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
	if len(segments) == 1 {
		if segments[0] == "KITCHEN" {
			return createKitchen(ws)
		}
	}
	if len(segments) != 2 {
		return nil, errors.New("Invalid payload format")
	}

	switch segments[0] {
	case "CUSTOMER":
		return createCustomer(segments[1], ws)
	case "WAITER":
		return createWaiter(segments[1], ws)
	default:
		return Customer{}, errors.New("Invalid identifier given")
	}
}

// createCustomer creates a customer with the given websocket with the given table number.
// An error is returned if a customer with the same table number already exists with an open connection.
func createCustomer(arg string, ws *websocket.Conn) (User, error) {
	n, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		return Customer{}, errors.New("Table number must be a number")
	}
	customers.Lock()
	defer customers.Unlock()
	for _, c := range customers.users {
		c := c.(Customer)
		if int(c.table) == int(n) {
			return Customer{}, fmt.Errorf("Table number %d is already connected", n)
		}
	}
	return Customer{ws: ws, table: uint(n)}, nil
}

// Remove removes this customer from the server register
func (c Customer) Remove() {
	customers.Lock()
	defer customers.Unlock()
	for i, cu := range customers.users {
		cu := cu.(Customer)
		if cu.table == c.table {
			// Remove the customer
			customers.users[i] = customers.users[len(customers.users)-1]
			customers.users = customers.users[:len(customers.users)-1]
			return
		}
	}
}

// createWaiter creates a waiter connection with the given websocket with the given name
func createWaiter(arg string, ws *websocket.Conn) (User, error) {
	if len(arg) == 0 {
		return Waiter{}, errors.New("Waiter name not provided")
	}
	waiters.Lock()
	defer waiters.Unlock()
	for _, w := range waiters.users {
		w := w.(Waiter)
		if w.name == arg {
			return Waiter{}, errors.New("Waiter with that name is already connected")
		}
	}
	return Waiter{ws: ws, name: arg}, nil
}

// Remove removes a waiter from the server register
func (w Waiter) Remove() {
	waiters.Lock()
	defer waiters.Unlock()
	for i, wa := range waiters.users {
		wa := wa.(Waiter)
		if wa.name == w.name {
			// Remove the waiter
			waiters.users[i] = waiters.users[len(waiters.users)-1]
			waiters.users = waiters.users[:len(waiters.users)-1]
			return
		}
	}
}

// createKitchen creates a kitchen from the given websocket and returns it
func createKitchen(ws *websocket.Conn) (User, error) {
	return Kitchen{ws: ws}, nil
}

// Remove removes a kitchen connection from the server register
func (k Kitchen) Remove() {
	kitchens.Lock()
	defer kitchens.Unlock()
	for i, ki := range kitchens.users {
		ki := ki.(Kitchen)
		if ki.ws == k.ws {
			// Remove the kitchen staff
			kitchens.users[i] = kitchens.users[len(kitchens.users)-1]
			kitchens.users = kitchens.users[:len(kitchens.users)-1]
			return
		}
	}
}

// Append appends to a thread-safe registry
func (r *Registry) Append(u User) {
	r.Lock()
	defer r.Unlock()
	r.users = append(r.users, u)
}
