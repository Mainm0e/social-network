package sockets

import "github.com/gorilla/websocket"

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
