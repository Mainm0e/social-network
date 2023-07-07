package sockets

import (
	"backend/db"
	"backend/events"
	"backend/handlers"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

/********************** PRIVATE MESSAGE LOGIC *******************************/

/*
BroadcastProvateMessage() is a method of the Manager struct, which takes a
receiver ID int and a json byte array as parameters. It then broadcasts the json
byte array to the the client with the given receiver ID. It returns an error
value, which is non-nil if any of the broadcasting operations failed or if
the receiver was not found.
*/
func (m *Manager) BroadcastPrivateMsg(receiverID int, msgEventJSON []byte) error {
	var sent bool
	// The range function on a sync.map accepts a function of the form
	// func(key, value interface{}) bool, which it calls once for each
	// item in the map. If the function returns false, the iteration stops.
	m.Clients.Range(func(key, client interface{}) bool {
		if client.(*Client).ID == receiverID {
			select {
			case client.(*Client).Egress <- msgEventJSON:
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

/********************** CHAT HISTORY / REQUEST LOGIC *************************/

/*
UnmarshalEventToChatHistoryRequest() takes an events.Event as input and
unmarshals it into a ChatHistoryRequest struct. It returns a pointer to the
ChatHistoryRequest struct and an error value, which is non-nil if the
unmarshalling operation failed.
*/
func UnmarshalEventToChatHistoryRequest(chatHistoryRequestEvent events.Event) (*ChatHistoryRequest, error) {
	var chatHistoryRequest ChatHistoryRequest
	if err := json.Unmarshal(chatHistoryRequestEvent.Payload, &chatHistoryRequest); err != nil {
		return nil, fmt.Errorf("UnmarshalEventToChatHistoryRequest() error: %v", err)
	}
	return &chatHistoryRequest, nil
}

/*
UnmarshalEventToChatHistory() takes an events.Event as input and unmarshals it
into a ChatHistory struct. It returns a pointer to the ChatHistory struct and
an error value, which is non-nil if the unmarshalling operation failed.
*/
func UnmarshalEventToChatHistory(chatHistoryEvent events.Event) (*ChatHistory, error) {
	var chatHistory ChatHistory
	if err := json.Unmarshal(chatHistoryEvent.Payload, &chatHistory); err != nil {
		return nil, fmt.Errorf("UnmarshalEventToChatHistory() error: %v", err)
	}
	return &chatHistory, nil
}

/*
FetchChatHistory() takes a ChatHistoryRequest struct as input and fetches the
chat history from the database. It returns a pointer to a ChatHistory struct
and an error value, which is non-nil if the database query failed. It works by
first preparing the condition and arguments for the SQL query based on the
ChatHistoryRequest struct (differentiating by "private" or "group" chat), then
fetching the chat history from the database, and finally casting the data into
a slice of db.Message structs.
*/
func FetchChatHistory(chatHistoryRequest ChatHistoryRequest) (*ChatHistory, error) {
	// Prepare the condition and arguments for the SQL query
	var condition string
	var args []any
	if chatHistoryRequest.ChatType == "private" {
		condition = "(senderId = ? AND receiverId = ?) OR (senderId = ? AND receiverId = ?)"
		args = []any{chatHistoryRequest.ClientID, chatHistoryRequest.TargetID, chatHistoryRequest.TargetID, chatHistoryRequest.ClientID}
	} else if chatHistoryRequest.ChatType == "group" {
		condition = "receiverId = ?"
		args = []any{chatHistoryRequest.TargetID}
	} else {
		return nil, fmt.Errorf("sockets.FetchChatHistory() error - unknown chat type: %s", chatHistoryRequest.ChatType)
	}

	// Fetch the chat history from the database
	data, err := db.FetchData("messages", condition, args...)
	if err != nil {
		return nil, fmt.Errorf("sockets.FetchChatHistory() error - %v", err)
	}

	// Cast the data into a slice of db.Message
	chatHistory := make([]db.Message, len(data))
	for i, item := range data {
		msg, ok := item.(db.Message)
		if !ok {
			return nil, fmt.Errorf("sockets.FetchChatHistory() error - Failed to cast item to Message")
		}
		chatHistory[i] = msg
	}

	// Compile the chat history into a ChatHistory struct and return it
	return &ChatHistory{
		ChatType:    chatHistoryRequest.ChatType,
		ClientID:    chatHistoryRequest.ClientID,
		TargetID:    chatHistoryRequest.TargetID,
		ChatHistory: chatHistory,
	}, nil
}

/*
SendChatHistory() is a method of the Manager struct that takes a pointer to a
ChatHistory struct as input and sends the chat history to the client with the
matching ID. It returns an error value, which is non-nil if the operation failed,
either in the wrapping of the ChatHistory struct into an events.Event struct,
the marshalling of the event into JSON, or the sending of the event to the client.
*/
func (m *Manager) SendChatHistory(chatHistory *ChatHistory) error {
	// Wrap the ChatHistory struct in an events.Event struct
	chatHistoryEvent, err := chatHistory.WrapEvent()
	if err != nil {
		return fmt.Errorf("sockets.SendChatHistory() error - %v", err)
	}

	// Marhsal the event into JSON
	eventJSON, err := json.Marshal(chatHistoryEvent)
	if err != nil {
		return fmt.Errorf("sockets.SendChatHistory() error - %v", err)
	}

	// Find the client with the matching ID and send the event
	m.Clients.Range(func(key interface{}, value interface{}) bool {
		client := value.(*Client)
		if client.ID == chatHistory.ClientID {
			select {
			case client.Egress <- eventJSON:
			default:
				// The Egress channel could be full or closed, or the client could be disconnected
				log.Printf("SendChatHistory() error - Could not send message to client %d\n", chatHistory.ClientID)
			}
			return false // Stop ranging as we found the client
		}
		return true // Continue ranging
	})

	return nil

}

/********************** COMMON LOGIC / FUNCTIONS *****************************/

func UnmarshalEventToChatMsg(msgEvent events.Event) (ChatMsg, error) {
	// Check the message type and unmarshal accordingly

	// Private message
	if msgEvent.Type == "privateMsg" {
		var pMsg PrivateMsg
		if err := json.Unmarshal(msgEvent.Payload, &pMsg); err != nil {
			return nil, fmt.Errorf("UnmarshalEventToChatMsg() error: %v", err)
		}
		return &pMsg, nil
	} else if msgEvent.Type == "groupMsg" {
		var gMsg GroupMsg
		if err := json.Unmarshal(msgEvent.Payload, &gMsg); err != nil {
			return nil, fmt.Errorf("UnmarshalEventToChatMsg() error: %v", err)
		}
		return &gMsg, nil
	}
	// Else, return an error
	return nil, fmt.Errorf("UnmarshalEventToChatMsg() error: invalid message type")
}

/*
RecordMsgToDB() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and records it to the database. It returns an error value, which is
non-nil if any of the database operations failed.
*/
func RecordMsgToDB(msg ChatMsg) error {
	msgType := msg.GetMsgType()
	senderID := msg.GetSenderID()
	receiverID := msg.GetReceiverID()
	messageContent := msg.GetMessage()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Insert message into database
	_, err := db.InsertData("messages", senderID, receiverID, messageContent, timestamp, msgType)
	if err != nil {
		return fmt.Errorf("RecordMsgToDB() error - unable to record message to database: %v", err)
	}

	return nil
}

/*
BroadcastMessage() takes a ChatMsg interface, which is either a PrivateMsg or a
GroupMsg, and broadcasts it to all clients in the chat. It returns an error value,
which is non-nil if any of the broadcasting operations failed.
*/
func (m *Manager) BroadcastMessage(msg ChatMsg) error {
	msgType := msg.GetMsgType()
	receiverID := msg.GetReceiverID()
	msgEvent := msg.WrapMsg()

	// Convert event to JSON
	msgEventJSON, err := json.Marshal(msgEvent)
	if err != nil {
		return err
	}

	// Decide broadcasting logic based on message type
	switch msgType {

	// For private messages, only broadcast to the receiver
	case "PrivateMsg":
		err = m.BroadcastPrivateMsg(receiverID, msgEventJSON)
		if err != nil {
			return fmt.Errorf("BroadcastMessage() error - %v", err)
		}

	// For group messages, broadcast to all clients within a group
	case "GroupMsg":
		err = m.BroadcastGroupMsg(receiverID, msgEventJSON)
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
func (m *Manager) HandleChatEvent(chatEvent events.Event, client *Client) error {
	// Unmarshal the event into a ChatMsg interface
	msg, err := UnmarshalEventToChatMsg(chatEvent)
	if err != nil {
		return fmt.Errorf("HandleChatMessage() error - %v", err)
	}

	// Store message in database
	err = RecordMsgToDB(msg)
	if err != nil {
		return fmt.Errorf("HandleChatMessage() error - %v", err)
	}

	// Broadcast message to clients
	err = m.BroadcastMessage(msg)
	if err != nil {
		return fmt.Errorf("HandleChatMessage() error - %v", err)
	}

	return nil
}

// TODO: HandleChatHistoryRequestEvent(): Receive request event and broadcast chat history

// TODO: HandleIsTypingEvent(): Receive isTyping event and broadcast to all clients in chat
