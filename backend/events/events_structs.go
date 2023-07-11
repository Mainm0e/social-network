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
