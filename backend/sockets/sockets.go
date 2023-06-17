package sockets

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

/*
NewManager is a function that creates a new Manager. This function doesn't take any
parameters and returns a pointer to a newly created Manager. It initializes all
the fields in the Manager struct - channels for broadcasting messages and
registering/unregistering clients, the client list, and the event handler map.
This function is typically called when the WebSocket server starts up and
needs to create a Manager to manage clients and messages.
*/
func NewManager() *Manager {
	return &Manager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    ClientList{},
		Handlers:   make(map[string]EventHandler),
	}
}

/*
NewClient is a function that creates a new Client. It takes two parameters:
a pointer to a websocket.Conn (which represents the WebSocket connection
between the server and the client) and a pointer to a Manager (which manages
the client and other clients). The function returns a pointer to a newly
created Client.  This function is typically called after a new WebSocket
connection has been established and a new Client needs to be created to
manage the connection.
*/
func NewClient(conn *websocket.Conn, wsManager *Manager) *Client {
	return &Client{
		Connection: conn,
		Manager:    wsManager,
		Egress:     make(chan []byte),
	}
}

/*
ReadData() is a method for a *Client struct, and starts a loop to continuously
read data from the client's websocket connection and react to that data.
*/
func (c *Client) ReadData() {
	defer func() {
		c.Manager.Unregister <- c
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(MAX_DATA_SIZE)

	c.Connection.SetReadDeadline(time.Now().Add(PONG_WAIT))

	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})

	// Infinite loop to continuously read data from the websocket connection
	for {
		// ReadMessage is a helper method that reads a message from the connection.
		// It returns the type of the message and the message itself, which is a byte slice.
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		// Unmarshal the received message into an event.
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error unmarshalling event: %v", err)

			// In case of an error unmarshalling, it sends back an error event to the client.
			errorEvent := Event{
				Type:    "", // TODO: Add error event type
				Payload: json.RawMessage(fmt.Sprintf(`{"error": "Failed to parse event: %v"}`, err)),
			}

			eventBytes, _ := json.Marshal(errorEvent)

			// Write the error event to the client's egress channel.
			c.Egress <- eventBytes
			continue
		}

		// If the event type exists in the map of handlers, execute the handler.
		handler, ok := c.Manager.Handlers[event.Type]
		if !ok {
			log.Printf("No handler for event type %v", event.Type)
			break
		}

		err = handler(event, c)
		if err != nil {
			log.Printf("Error handling event: %v", err)
			break
		}
	}
}
