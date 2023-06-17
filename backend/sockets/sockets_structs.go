package sockets

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

/*
Manager struct is a central place that keeps track of all connected clients, broadcasts
messages to all clients, and handles the registration and unregistration of clients.
It includes a channel for broadcast data, channels for client registration and
unregistration, a synchronized map of connected clients, a map of event handlers.
*/
type Manager struct {
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Clients    ClientList
	Handlers   map[string]EventHandler
}

/*
An upgrader object to upgrade HTTP connections to WebSocket connections. This
allows the server to respond to HTTP upgrade requests from clients that want to
initiate a WebSocket connection.
*/
var websocketUpgrader = websocket.Upgrader{
	// Allow connections from any origin
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*
Event struct represents a single event that can occur within a WebSocket
communication, such as a user sending or receiving a message.
*/
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

/*
EventHandler type is a function that handles an event. It takes an Event and a Client
as parameters and returns an error if the event cannot be handled successfully.
*/
type EventHandler func(event Event, client *Client) error

/*
ClientList is a struct that represents a map of connected clients. It is a wrapper
around a sync.Map, which is a built-in type that provides a concurrent-safe map and
thus can be used by multiple goroutines without additional locking or coordination.
It is optimized for use cases in which keys and values are only ever added to the map,
retrieved, and deleted.
*/
type ClientList struct {
	sync.Map
}

/*
Client struct represents a single client that is connected to the server via a
WebSocket connection. It includes the client's WebSocket connection, the Manager
that manages the client, and a channel for outgoing messages (egress).
*/
type Client struct {
	Connection *websocket.Conn
	Manager    *Manager
	Egress     chan []byte // A channel for outgoing messages
}
