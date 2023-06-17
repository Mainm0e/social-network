/*
Package "sockets" provides broad functionality for managing WebSocket connections, and
its use can be divided as follows:

 1. Initializing the WebSocket Manager:
    The first step will be initializing the WebSocket Manager, which is a central place
    that keeps track of all connected clients, broadcasts messages to all clients, and
    handles the registration and unregistration of clients. This is done by calling the
    NewManager() function and should be done when a server starts up.

 2. Registering Event Handlers:
    Next, event handlers need to be registered. These handlers are functions that define
    what to do when a certain type of event is received. Events could be anything one
    defines, such as a new chat message being sent, a new user joining the forum, or a
    user leaving the forum. These handlers should be defined separately and then
    registered with the Manager by calling RegisterHandler(). This should also be done during
    server startup, after the Manager is initialized.

 3. Running the Manager:
    After the Manager is initialized and event handlers are registered, the Manager can start
    its main loop by calling Run. This is typically done in a new goroutine, as it's a
    blocking function that runs indefinitely.

 4. Serving WebSocket Connections:
    ServeWS() is an HTTP handler that upgrades an incoming HTTP(S) request to a WebSocket
    connection, and needs to be attached to a route on ones HTTP(S) server. When a client
    makes a request to this route, ServeWS will upgrade the connection and register the
    client with the Manager. After a client is registered, ReadData and WriteData routines
    will start for this client in separate goroutines.

 5. Sending and Receiving Data:
    The ReadData and WriteData methods on the Client struct handle reading and writing
    WebSocket messages from/to the client. ReadData reads messages from the WebSocket
    connection, unmarshals the JSON into an Event struct, and then calls the appropriate
    handler based on the event type. WriteData writes messages to the WebSocket connection.
    These messages could be sent by the server or could be broadcast messages from other
    clients.

The "sockets" package is fairly self-contained within the greater project go-module. Other parts
of ones server don't necessarily need to know about it, except for registering the ServeWS handler
with a route on the server. However, one needs to write the event handlers that get registered
with the Manager, and those handlers will likely interact with other parts of ones server. For
example, a chat message handler might need to interact with the database package to store the
message. This means that certain aspects should be defined in a common package, such as the Event
struct, which is used by both the "sockets" package and packages used for handlers etc. The same
goes for actual event definitions, such as a chat message event, and should be defined in a common
package.

The NewClient, ReadData, and WriteData methods are all used internally within the "sockets"
package and don't need to be called directly from outside the package.
*/
package sockets
