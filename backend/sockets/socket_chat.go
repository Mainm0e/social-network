package sockets

import (
	"encoding/json"
	"errors"
)

/********************** PRIVATE MESSAGE LOGIC *******************************/

/*
UnmarshalJSONToPrivateMsg() takes a json byte array, which is usually received
from the frontend in the form of a websocket message, and unmarshals it into
a PrivateMsg struct. It then returns a pointer to this PrivateMsg struct, along
with an error value, which is non-nil if the unmarshalling process failed.
*/
func UnmarshalJSONToPrivateMsg(jsonMsg []byte) (*PrivateMsg, error) {
	var privateMsg PrivateMsg
	if err := json.Unmarshal(jsonMsg, &privateMsg); err != nil {
		return nil, err
	}
	return &privateMsg, nil
}

/*
BroadcastProvateMessage() is a method of the Manager struct, which takes a
receiver ID and a json byte array as parameters. It then broadcasts the json
byte array to the the client with the given receiver ID. It returns an error
value, which is non-nil if any of the broadcasting operations failed or if
the receiver was not found.
*/
func (m *Manager) BroadcastPrivateMsg(receiver string, payloadJSON []byte) error {
	var sent bool
	// The range function on a sync.map accepts a function of the form
	// func(key, value interface{}) bool, which it calls once for each
	// item in the map. If the function returns false, the iteration stops.
	m.Clients.Range(func(key, client interface{}) bool {
		if client.(*Client).ID == receiver {
			select {
			case client.(*Client).Egress <- payloadJSON:
				sent = true
			default:
				close(client.(*Client).Egress)
				m.Clients.Delete(key)
				return false
			}
			return false // Stop iteration after the message has been sent to target client
		}
		return true
	})

	if !sent {
		return errors.New("message could not be sent: receiver not found")
	}

	return nil
}

/********************** GROUP MESSAGE LOGIC **********************************/

/********************** COMMON LOGIC / FUNCTIONS *****************************/

/*
RecordMsgToDB() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and records it to the database. It returns an error value, which is
non-nil if any of the database operations failed.
*/
func RecordMsgToDB(msg ChatMsg) error {
	sender := msg.GetSender()
	receiver := msg.GetReceiver()
	message := msg.GetMessage()
	timestamp := msg.GetTimestamp()

	// TODO: Write database logic to record message to database

	return nil
}

/*
BroadcastMessage() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and broadcasts it to all clients in the chat. It returns an error value,
which is non-nil if any of the broadcasting operations failed.
*/
func (m *Manager) BroadcastMessage(msg ChatMsg) error {
	return nil
}

/*
HandleMessage() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and handles it. This means that it records the message to the database
and broadcasts it to all clients in the chat. It returns an error value, which
is non-nil if any of the operations failed.
*/
func (m *Manager) HandleMessage(msg ChatMsg) error {
	// Store message in database
	err := RecordMsgToDB(msg)
	if err != nil {
		return err
	}

	// Broadcast message to clients
	err = m.BroadcastMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
