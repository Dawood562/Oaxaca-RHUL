package endpoints

import (
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

// NewConnection registers a new connection to the system given its websocket connection
func NewConnection(ws *websocket.Conn) error {
	_, m, err := ws.ReadMessage()
	if err != nil {
		return err
	}
	t := getTableNum(string(m))
	if t == -1 {
		return ws.WriteMessage(websocket.TextMessage, []byte("DENIED"))
	}
	return ws.WriteMessage(websocket.TextMessage, []byte("WELCOME"))
}

// getTableNum extracts a table number from a message or returns -1 if the message is invalid
func getTableNum(m string) int {
	segments := strings.Split(string(m), ":")
	if len(segments) != 2 {
		return -1
	}
	if segments[0] != "CUSTOMER" {
		return -1
	}
	// Extract the table number
	ret, err := strconv.ParseInt(segments[1], 10, 32)
	if err != nil {
		return -1
	}
	return int(ret)
}
