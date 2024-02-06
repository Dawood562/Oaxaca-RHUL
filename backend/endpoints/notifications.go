package endpoints

import (
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

type Customer struct {
	ws    *websocket.Conn
	table uint
}

type Waiter struct {
	ws   *websocket.Conn
	name string
}

var customers []Customer
var waiters []Waiter

// NewConnection registers a new connection to the system given its websocket connection
func NewConnection(ws *websocket.Conn) error {
	_, m, err := ws.ReadMessage()
	if err != nil {
		return err
	}
	r := registerConnection(string(m), ws)
	if !r {
		return ws.WriteMessage(websocket.TextMessage, []byte("DENIED"))
	}
	return ws.WriteMessage(websocket.TextMessage, []byte("WELCOME"))
}

// registerConnection returns true if the initial message from the connection is OK
func registerConnection(m string, ws *websocket.Conn) bool {
	segments := strings.Split(string(m), ":")
	if len(segments) != 2 {
		return false
	}

	switch segments[0] {
	case "CUSTOMER":
		// Extract the table number
		ret, err := strconv.ParseInt(segments[1], 10, 32)
		if err != nil {
			return false
		}
		customers = append(customers, Customer{ws: ws, table: uint(ret)})
	case "WAITER":
		if len(segments[1]) == 0 {
			return false
		}
		waiters = append(waiters, Waiter{ws: ws, name: segments[1]})
	}

	return true
}
