package sockets

import "encoding/json"

/* COMMON LOGIC / FUNCTIONS */

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

/* PRIVATE MESSAGE LOGIC */

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
