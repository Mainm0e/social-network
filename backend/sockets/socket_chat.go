package sockets

import (
	"backend/handlers"
	"encoding/json"
	"errors"
	"fmt"
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
receiver ID int and a json byte array as parameters. It then broadcasts the json
byte array to the the client with the given receiver ID. It returns an error
value, which is non-nil if any of the broadcasting operations failed or if
the receiver was not found.
*/
func (m *Manager) BroadcastPrivateMsg(receiverID int, payloadJSON []byte) error {
	var sent bool
	// The range function on a sync.map accepts a function of the form
	// func(key, value interface{}) bool, which it calls once for each
	// item in the map. If the function returns false, the iteration stops.
	m.Clients.Range(func(key, client interface{}) bool {
		if client.(*Client).ID == receiverID {
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

/*
UmarshalJSONToGroupMsg() takes a json byte array, which is usually received
from the frontend in the form of a websocket message, and unmarshals it into
a GroupMsg struct. It then returns a pointer to this GroupMsg struct, along
with an error value, which is non-nil if the unmarshalling process failed.
*/
func UnmarshalJSONToGroupMsg(jsonMsg []byte) (*GroupMsg, error) {
	var groupMsg GroupMsg
	if err := json.Unmarshal(jsonMsg, &groupMsg); err != nil {
		return nil, err
	}
	return &groupMsg, nil
}

/*
BroadcastGroupMessage() is a method of the Manager struct, which takes a
groupID int and a json byte array as input parameters. It then broadcasts
the json byte array to all clients in the group chat. It returns an error
value, which is non-nil if any of the broadcasting operations failed or if
there are no members in the group.
*/
func (m *Manager) BroadcastGroupMsg(groupID int, payloadJSON []byte) error {
	// Retrieve the userIDs of the group members from the database
	memberUserIDs, err := handlers.GetAllGroupMemberIDs(groupID)
	if err != nil {
		return fmt.Errorf("BroadCastGroupMsg() error - unable to retrieve group "+
			"\" %v \" members: %v", groupID, err)
	}

	// Check if there are any members in the group
	if len(memberUserIDs) == 0 {
		return fmt.Errorf("BroadCastGroupMsg() - no members in group "+
			"\" %v \": message could not be broadcast", groupID)
	}

	// Flag to track if message was sent to at least one member
	var sent bool = false

	// Loop through the memberUserIDs
	for _, userID := range memberUserIDs {
		// For group messages the function always returns true to keep the
		// iteration going and send messages to all clients in the group chat.
		m.Clients.Range(func(key, client interface{}) bool {
			// Check if this client's ID matches the current userID
			if client.(*Client).ID == userID {
				// Attempt to send the message to this client
				select {
				case client.(*Client).Egress <- payloadJSON:
					sent = true
				default:
					close(client.(*Client).Egress)
					m.Clients.Delete(key)
				}
				// Stop iteration for this client as the message has been sent (or attempted)
				return false
			}
			// Continue iteration for other clients
			return true
		})
	}

	// Check if the message was sent to at least one member
	if !sent {
		return fmt.Errorf("BroadCastGroupMsg() - no active connections in group "+
			"\" %v \": message could not be broadcast", groupID)
	}

	// Return nil if there were no errors
	return nil
}

/********************** COMMON LOGIC / FUNCTIONS *****************************/

/*
RecordMsgToDB() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and records it to the database. It returns an error value, which is
non-nil if any of the database operations failed.
*/
func RecordMsgToDB(msg ChatMsg) error {
	msgType := msg.GetMsgType()
	senderID := msg.GetSenderID()
	receiverID := msg.GetReceiverID()
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
	senderID := msg.GetSenderID()
	receiverID := msg.GetReceiverID() // Can be a userID or groupID
	message := msg.GetMessage()
	timestamp := msg.GetTimestamp()
	msgType := msg.GetMsgType()

	// payload to be sent to clients
	payload := map[string]interface{}{
		"type": msgType,
		"payload": map[string]string{
			"senderID":   string(senderID),
			"receiverID": string(receiverID),
			"message":    message,
			"timeStamp":  timestamp,
		},
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Decide broadcasting logic based on message type
	switch msgType {

	// For private messages, only broadcast to the receiver
	case "PrivateMsg":
		err = m.BroadcastPrivateMsg(receiverID, payloadJSON)
		if err != nil {
			return fmt.Errorf("BroadcastMessage() error - %v", err)
		}

	// For group messages, broadcast to all clients within a group
	case "GroupMsg":
		err = m.BroadcastGroupMsg(receiverID, payloadJSON)
		if err != nil {
			return fmt.Errorf("BroadcastMessage() error - %v", err)
		}
	}

	return nil
}

/*
HandleChatMessage() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and handles it. This means that it records the message to the database
and broadcasts it to all clients in the chat. It returns an error value, which
is non-nil if any of the operations failed.
*/
func (m *Manager) HandleChatMessage(msg ChatMsg) error {
	// Store message in database
	err := RecordMsgToDB(msg)
	if err != nil {
		return fmt.Errorf("HandleMessage() error - %v", err)
	}

	// Broadcast message to clients
	err = m.BroadcastMessage(msg)
	if err != nil {
		return fmt.Errorf("HandleMessage() error - %v", err)
	}

	return nil
}
