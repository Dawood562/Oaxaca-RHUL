//go:build integration

package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/contrib/websocket"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestOpenWebsocket(t *testing.T) {
	// Mockup the app
	var mockFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		upgrader := gwebsocket.Upgrader{}
		ws, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(t, err)
		defer ws.Close()

		// Send test debug message
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome, Table")))
	}
	s := httptest.NewServer(mockFunc)
	defer s.Close()
	url, err := url.Parse(s.URL)
	assert.NoError(t, err)
	url.Scheme = "ws"

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
			ws, _, err := gwebsocket.DefaultDialer.Dial(url.String()+"/notifications", nil)
			assert.NoError(t, err, "Test that creating a websocket connection does not create an error")
			defer ws.Close()
			// Try reading from the websocket
			_, _, err = ws.ReadMessage()
			assert.NoError(t, err, "Test that receiving a message from the websocket does not create an error")
		})
	}
}
