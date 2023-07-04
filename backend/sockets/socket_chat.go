package sockets

import "encoding/json"

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
	receiver := msg.GetReceiver()
	sender := msg.GetSender()
	message := msg.GetMessage()
	timestamp := msg.GetTimestamp()
	msgType := msg.GetType()

	// payload to be sent to clients
	payload := map[string]interface{}{
		"type": msgType,
		"payload": map[string]string{
			"senderUsername": sender,
			"receiver":       receiver,
			"message":        message,
			"timeStamp":      timestamp,
		},
	}

	// TODO: Logic to broadcast message to all clients in the chat

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
