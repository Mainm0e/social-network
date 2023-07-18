package events

import "encoding/json"

/*
Event struct represents a single event that can occur within a WebSocket
or HTTP/S communication. It includes a type and a payload.
*/
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

/***************************** EVENT PAYLOAD STRUCTS *****************************/

/*
ErrorMessage struct represents an error message that can be sent to the
frontend / client. It includes a message and a status code.
*/
type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
