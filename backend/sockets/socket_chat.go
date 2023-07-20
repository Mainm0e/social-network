package sockets

import (
	"backend/db"
	"backend/events"
	"backend/handlers"
	"backend/server/sessions"
	"encoding/json"
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
func (m *Manager) BroadcastPrivateMsg(senderID int, receiverID int, msgEventJSON []byte) error {
	// Flags to check if sender and receiver have been updated
	senderUpdated := false
	receiverUpdated := false

	// Iterate through the map of clients
	m.Clients.Range(func(key, client interface{}) bool {
		// Check if the current client is the sender or the receiver
		if client.(*Client).ID == senderID || client.(*Client).ID == receiverID {
			// Send the message to the client
			client.(*Client).Egress <- msgEventJSON
			// Set the respective flag to true
			if client.(*Client).ID == senderID {
				log.Printf("BroadcastPrivateMsg() to sender with ID: \" %v \" ", senderID)
				senderUpdated = true
			} else if client.(*Client).ID == receiverID {
				log.Printf("BroadcastPrivateMsg() to receiver with ID: \" %v \" ", receiverID)
				receiverUpdated = true
			}
			// Do not stop iteration; there may be another client to send to
		}
		return true // Continue iteration regardless
	})

	// Check if both the sender and receiver have received the message
	if !senderUpdated && !receiverUpdated {
		return fmt.Errorf("BroadcastPrivateMsg() error - sender \" %v \" and receiver \" %v \" not updated as BOTH respective live socket connections not found", senderID, receiverID)
	} else if !senderUpdated {
		log.Printf("BroadcastPrivateMsg() notice - sender with ID: \" %v \" not updated in private chat with client \" %v \" as live socket connection not found", senderID, receiverID)
	} else if !receiverUpdated {
		log.Printf("BroadcastPrivateMsg() notice - receiver with ID: \" %v \" not updated in private chat with client \" %v \" as live socket connection not found", receiverID, senderID)
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
					log.Printf("BroadcastGroupMsg() to sender with ID: \" %v \" ", userID)
					sent = true
				default:
					close(client.(*Client).Egress)
					m.Clients.Delete(key)
					log.Printf("BroadcastGroupMsg() - deleting client with ID: \" %v \" ", userID)
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
func UnmarshalEventToChatHistoryRequest(chatHistoryRequestEvent events.Event) (ChatHistoryRequest, error) {
	var chatHistoryRequest ChatHistoryRequest
	if err := json.Unmarshal(chatHistoryRequestEvent.Payload, &chatHistoryRequest); err != nil {
		return chatHistoryRequest, fmt.Errorf("UnmarshalEventToChatHistoryRequest() error: %v", err)
	}
	return chatHistoryRequest, nil
}

/*
UnmarshalEventToChatHistory() takes an events.Event as input and unmarshals it
into a ChatHistory struct. It returns a pointer to the ChatHistory struct and
an error value, which is non-nil if the unmarshalling operation failed.

**NOTE** This function is currently not used anywhere, but it is kept here
for future use.
*/
func UnmarshalEventToChatHistory(chatHistoryEvent events.Event) (ChatHistory, error) {
	var chatHistory ChatHistory
	if err := json.Unmarshal(chatHistoryEvent.Payload, &chatHistory); err != nil {
		return chatHistory, fmt.Errorf("UnmarshalEventToChatHistory() error: %v", err)
	}
	return chatHistory, nil
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
	clientUpdated := false
	m.Clients.Range(func(key interface{}, value interface{}) bool {
		client, ok := value.(*Client)
		if !ok {
			log.Printf("SendChatHistory() error - Could not cast value to *Client")
			// Handle the error if need be... The value is not of type *Client.
			return true
		}
		if client.ID == chatHistory.ClientID {
			client.Egress <- eventJSON
			log.Printf("SendChatHistory() - Sent chat history to client %d\n", chatHistory.ClientID)
			clientUpdated = true
			return false // Stop ranging as the client is found
		}
		return true // Continue ranging
	})

	if !clientUpdated {
		return fmt.Errorf("SendChatHistory() error - client with ID: \" %v \" not found", chatHistory.ClientID)
	}

	return nil
}

/********************** IS TYPING EVENT / LOGIC ******************************/

/*
HandleIsTypingEvent() is a method of the Manager struct that takes an events.Event
and a pointer to a Client as input. It unmarshals the event payload into an
IsTyping struct, wraps the event in a JSON payload, and sends it to the relevant
client(s) (i.e. the client(s) in the same chat as the client that sent the event).
*/
func (m *Manager) HandleIsTypingEvent(typingEvent events.Event, client *Client) {
	// Initialise an IsTyping struct and unmarshal the event payload into it
	var isTyping IsTyping
	if err := json.Unmarshal(typingEvent.Payload, &isTyping); err != nil {
		log.Printf("HandleIsTypingEvent() - Error unmarshalling event payload: %v", err)
		return
	}

	// Marshal the event into JSON bytes
	eventBytes, err := json.Marshal(typingEvent)
	if err != nil {
		log.Printf("HandleIsTypingEvent() - Error marshalling event: %v", err)
		return
	}

	// Notify the relevant client(s) about the typing status.
	m.Clients.Range(func(k, v interface{}) bool {
		otherClient, ok := v.(*Client)
		if !ok {
			log.Printf("HandleIsTypingEvent() - Error casting client from client list: %v", v.(*Client))
			return true // Continue iteration
		}

		if isTyping.ChatType == "private" && otherClient.ID == isTyping.TargetID {
			otherClient.Egress <- eventBytes
		} else if isTyping.ChatType == "group" {
			memberUserIDs, err := handlers.GetAllGroupMemberIDs(isTyping.TargetID)
			if err != nil {
				log.Printf("HandleIsTypingEvent() - Error getting group member IDs: %v", err)
				return true // Continue iteration
			}
			// Check if the otherClient.ID is in the list of member IDs
			for _, id := range memberUserIDs {
				if otherClient.ID == id {
					otherClient.Egress <- eventBytes
					break
				}
			}
		}
		return true // Continue iteration
	})
}

/********************** LOGOUT EVENT / LOGIC *********************************/

/*
HandleLogoutEvent() documentation...
*/
func (m *Manager) HandleLogoutEvent(logoutEvent events.Event, client *Client) {
	// Initialise a LogoutWS struct and unmarshal the event payload into it
	var logoutData handlers.UserCredential // TODO may need to change this to socket-specific struct
	if err := json.Unmarshal(logoutEvent.Payload, &logoutData); err != nil {
		log.Printf("HandleLogoutEvent() - Error unmarshalling event payload: %v", err)
		return
	}

	// Log the client out of their session
	err := sessions.Logout(logoutData.SessionId)
	if err != nil {
		log.Println("Error logging out:", err)
		return
	}

	// Send a logout confirmation to the client
	// TODO: This is not working for some reason (tried various versions)

	// Close the client's Egress channel and remove them from the active clients
	close(client.Egress)
	m.Clients.Delete(client.ID)
	log.Printf("HandleLogoutEvent() - Removed client %d from active clients\n", client.ID)
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
	// Log the broadcast attempt
	log.Printf("BroadcastMessage() - Broadcasting message: %v\n", msg)

	msgType := msg.GetMsgType()
	senderID := msg.GetSenderID()
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
		err = m.BroadcastPrivateMsg(senderID, receiverID, msgEventJSON)
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
and broadcasts it to all clients in the chat. It does not return anything, but
logs errors if any are encountered.
*/
func (m *Manager) HandleChatEvent(chatEvent events.Event, client *Client) {
	// Log the attempt to handle the event
	log.Printf("HandleChatMessage() - Received event: %v\n", chatEvent)

	// Unmarshal the event into a ChatMsg interface
	msg, err := UnmarshalEventToChatMsg(chatEvent)
	if err != nil {
		log.Printf("HandleChatMessage() error - %v", err)
		// TODO: Send error message to client
		return
	}

	// Store message in database
	err = RecordMsgToDB(msg)
	if err != nil {
		log.Printf("HandleChatMessage() error - %v", err)
		// TODO: Send error message to client
		return
	}

	// Broadcast message to clients
	err = m.BroadcastMessage(msg)
	if err != nil {
		log.Printf("HandleChatMessage() error - %v", err)
		// TODO: Send error message to client
		return
	}
}

/*
HandleChatHistoryRequestEvent() takes a ChatHistoryRequest event as input, along
with a pointer to the client that sent the request, and handles it. This means that
it first unmashals the event into a ChatHistoryRequest struct, then fetches the
chat history from the database, and finally sends the chat history to the client.
It does not return anything, but logs errors if any are encountered.
*/
func (m *Manager) HandleChatHistoryRequestEvent(chatHistoryRequestEvent events.Event, client *Client) {
	// Log the attempt to handle the event
	log.Printf("HandleChatHistoryRequestEvent() - Received event: %v\n", chatHistoryRequestEvent)

	// Unmarshal the event into a ChatHistoryRequest struct
	chatHistoryRequest, err := UnmarshalEventToChatHistoryRequest(chatHistoryRequestEvent)
	if err != nil {
		log.Printf("HandleChatHistoryRequestEvent() error - %v", err)
		// TODO: Send error message to client
		return
	}

	// Get chat history from database, returned as a pointer to a ChatHistory struct
	chatHistory, err := FetchChatHistory(chatHistoryRequest)
	if err != nil {
		log.Printf("HandleChatHistoryRequestEvent() error - %v", err)
		// TODO: Send error message to client
		return
	}

	// Send chat history to client
	err = m.SendChatHistory(chatHistory)
	if err != nil {
		log.Printf("HandleChatHistoryRequestEvent() error - %v", err)
		// TODO: Send error message to client
		return
	}
}

/*
SendErrorMessageToClient() is a method of the Manager struct. It takes a pointer
to a client, an error message string and a status code integer as input, and sends
the error message to the client. It returns nothing.
*/
func (m *Manager) SendErrorMessageToClient(client *Client, message string, statusCode int) {
	// Create ErrorMessage
	errMsg := events.ErrorMessage{
		Message:    message,
		StatusCode: statusCode,
	}

	// Package ErrorMessage into an Event
	event, err := events.PackageErrorEvent(errMsg)
	if err != nil {
		log.Printf("SendErrorMessageToClient error packaging error: %v", err)
		return
	}

	// Convert Event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("SendErrorMessageToClient error marshalling event: %v", err)
		return
	}

	// Send the Event JSON to client
	select {
	case client.Egress <- eventJSON:
	default:
		// The Egress channel could be full or closed, or the client could be disconnected
		log.Printf("SendErrorMessageToClient() error - Could not send error message to client %d\n", client.ID)
		return
	}
}
