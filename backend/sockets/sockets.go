package sockets

import (
	"backend/events"
	"backend/server/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	m := &Manager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    ClientList{},
		Handlers:   make(map[string]EventHandler),
	}

	// Add the chat handler to the Handlers map
	m.Handlers["privateMsg"] = m.HandleChatEvent
	m.Handlers["groupMsg"] = m.HandleChatEvent

	// Add the chat history handler to the Handlers map
	m.Handlers["chatHistoryRequest"] = m.HandleChatHistoryRequestEvent

	// Add the isTyping handler to the Handlers map
	m.Handlers["isTyping"] = m.HandleIsTypingEvent

	return m
}

/*
NewClient is a function that creates a new Client. It takes three parameters:
a pointer to a websocket.Conn (which represents the WebSocket connection between
the server and the client), a pointer to a Manager (which manages the client and
other clients), and a string (which is the client's ID). The function returns a
pointer to a newly created Client.  This function is typically called after a new
WebSocket connection has been established and a new Client needs to be created to
manage the connection.
*/
func NewClient(conn *websocket.Conn, wsManager *Manager, id int) *Client {
	return &Client{
		Connection: conn,
		Manager:    wsManager,
		Egress:     make(chan []byte),
		ID:         id,
	}
}

/*
ReadData() is a method for a *Client struct, and starts a loop to continuously
read data from the client's websocket connection and react to that data.
*/
func (c *Client) ReadData() {
	// Defer the closing of the client's websocket connection, which gets called
	// when the function returns
	defer func() {
		log.Printf("sockets.ReadData() - Closing websocket connection for client \" %v \"", c.ID)
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
				log.Printf("sockets.ReadData() - Unexpected close error: %v", err)
			}
			break
		}

		// Unmarshal the received message into an event.
		var event events.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("sockets.ReadData() - Error unmarshalling event: %v", err)

			// In case of an error unmarshalling, it sends back an error event to the client.
			errorEvent := events.Event{
				Type:    "error",
				Payload: json.RawMessage(fmt.Sprintf(`{"sockets.ReadData() error": "Failed to parse event: %v"}`, err)),
			}

			eventBytes, _ := json.Marshal(errorEvent)

			// Write the error event to the client's egress channel.
			c.Egress <- eventBytes
			continue
		}

		// If the event type exists in the map of handlers, execute the handler.
		handler, ok := c.Manager.Handlers[event.Type]
		if !ok {
			log.Printf("sockets.ReadData() - No handler for event type %v", event.Type)
			break
		}

		handler(event, c)
	}
}

/*
WriteData() is a method for a *Client struct, which starts a loop to
continuously write data to the client's websocket connection. It also
sends periodic pings to the client.
*/
func (c *Client) WriteData() {
	// Ticker is a timer that goes off (ticks) at regular intervals.
	ticker := time.NewTicker(PING_INTERVAL)

	// Defer the stopping of the ticker and closing of the client's websocket connection,
	// which gets called when the function returns.
	defer func() {
		log.Printf("sockets.WriteData() - Closing websocket connection for client \" %v \"", c.ID)
		ticker.Stop()
		c.Connection.Close()
	}()

	// Infinite loop to continuously write data to the websocket connection.
	for {
		// The select statement lets a goroutine wait on multiple communication
		// operations (channels). A select blocks until one of its cases can run,
		// then it executes that case. It chooses one at random if multiple are ready.
		select {

		// This case handles outgoing messages from the client to the websocket connection.
		case data, ok := <-c.Egress:
			// If the channel is closed, the ok variable will be set to false.
			if !ok {
				log.Printf("sockets.WriteData() - Egress channel unavailable for client \" %v \", websocket closed", c.ID)
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// NextWriter returns a writer for the next message to send.
			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// The data from the channel is written to the connection.
			w.Write(data)

			// If there are more messages in the channel, they are written to the connection as well.
			// This helps in flushing any queued messages in the channel.
			n := len(c.Egress)
			for i := 0; i < n; i++ {
				w.Write(<-c.Egress)
			}

			// Close finalizes the message. The writer must be closed before the next call to
			// NextWriter. Close returns an error if the message was not correctly formed.
			if err := w.Close(); err != nil {
				return
			}

		// This case handles the periodic sending of PingMessage to the client over the ws-connection.
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/*
Run is the main loop for the Manager. It listens for incoming actions
such as client registrations, unregistrations, and broadcasting messages.
*/
func (m *Manager) Run() {
	log.Println("sockets.Run() - Starting websocket manager")
	for {
		select {
		// A new client is registering: Store it in the clients map.
		case client := <-m.Register:
			log.Println("sockets.Run() - Registering new client")
			m.Clients.Store(client.ID, client)

		// A client is unregistering: If it exists in the clients map, remove it.
		case client := <-m.Unregister:
			log.Println("sockets.Run() - Deregistering new client")
			if _, ok := m.Clients.Load(client.ID); ok {
				m.Clients.Delete(client.ID)
				close(client.Egress)
			}

		// Data is being broadcast: Send it to all connected clients.
		case data := <-m.Broadcast:
			log.Println("sockets.Run() - Broadcasting data")
			m.Clients.Range(func(key, value interface{}) bool {
				client := key.(*Client)
				select {
				case client.Egress <- data:
					// The data was sent successfully, continue to the next client.
					return true
				default:
					// The client's send channel is unavailable. Remove it.
					close(client.Egress)
					m.Clients.Delete(client.ID)
					return false // Stop iteration.
				}
			})
		}
	}
}

/*
ServeWS is an HTTP handler function that upgrades an HTTP(S) connection to
a WebSocket connection. It creates a new client and then initiates the
reading and writing goroutines for that client. Parameters:
- w: The HTTP ResponseWriter that the handler will use to send HTTP responses.
- r: The HTTP Request that has been received by the handler.
*/
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket initialisation started...")

	// Perform validation checks on the session cookie.
	cookie, err := r.Cookie(sessions.COOKIE_NAME)
	if err != nil {
		log.Printf("sessions.ServeWS() error - No sessionID cookie found: %v", err)
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	} else {
		isValid, err := sessions.CookieCheck(cookie)
		if !isValid || err != nil {
			log.Printf("sessions.ServeWS() error - Invalid sessionID cookie for session \" %v \": %v", cookie.Value, err)
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
	}
	sessionID := cookie.Value

	// Retrieve the associated UserID from the sessions Store.
	userSession, sessionExists, err := sessions.Store.Get(sessionID)
	if err != nil || !sessionExists {
		log.Printf("sockets.ServeWS() error - Failed to retrieve UserID for session \" %v \" from sessions Store: %v", sessionID, err)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	// Create a new client for the WebSocket connection.
	client := NewClient(conn, m, userSession.UserID)

	// Register the new client with the Manager.
	m.Register <- client

	// Starts the read and write goroutines for the client.
	go client.ReadData()
	go client.WriteData()

	log.Println("Websocket initialisation complete")
}
