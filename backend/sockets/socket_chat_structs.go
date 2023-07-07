package sockets

import (
	"backend/db"
	"time"
)

/************************** COMMON EVENT STRUCTS **************************/
type PrivateMsg struct {
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Message    string `json:"message"`
	Timestamp  string `json:"timeStamp"` // Received as an empty string from frontend
}

type GroupMsg struct {
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Message    string `json:"message"`
	Timestamp  string `json:"timeStamp"` // Received as an empty string from frontend
}

type IsTyping struct {
	ChatType string `json:"chatType"`
	ClientID int    `json:"clientID"` // UserID of the client that is typing
	TargetID int    `json:"targetID"` // UserID / GroupID to identify chat in which typing is happening
}

/*********************** ONLY FROM FRONTEND EVENT STRUCT ******************/

type ChatHistoryRequest struct {
	ChatType string `json:"chatType"`
	ClientID int    `json:"clientID"`
	TargetID int    `json:"targetID"` // UserID or GroupID
}

/*********************** ONLY FROM BACKEND EVENT STRUCT ******************/

type ChatHistory struct {
	ChatType    string       `json:"chatType"`
	ClientID    int          `json:"clientID"`
	TargetID    int          `json:"targetID"` // UserID or GroupID
	ChatHistory []db.Message `json:"chatHistory"`
}

/*********************** INTERFACES & METHODS ****************************/
/********* NOTE TO SELF: NEVER FORGET HOW COOL INTERFACES ARE! ***********/

// ChatMsg is an interface that is implemented by both PrivateMsg and GroupMsg.
// This allows us to use the same functions for both types of messages.
type ChatMsg interface {
	GetSenderID() int
	GetReceiverID() int
	GetMessage() string
	GetTimestamp() string
	GetMsgType() string

	// Define fields that are common to both PrivateMsg and GroupMsg
	WrapMsg() struct {
		Type    string `json:"type"`
		Payload struct {
			SenderID   int    `json:"senderID"`
			ReceiverID int    `json:"receiverID"`
			Message    string `json:"message"`
			Timestamp  string `json:"timeStamp"`
		} `json:"payload"`
	}
}

/*
GetSenderID() returns the senderID of the message. This is a common field for
both PrivateMsg and GroupMsg. This is useful when we want to send a message to
the frontend, but we don't know if it is a PrivateMsg or GroupMsg.
*/
func (p *PrivateMsg) GetSenderID() int {
	return p.SenderID
}

func (g *GroupMsg) GetSenderID() int {
	return g.SenderID
}

/*
GetReceiverID() returns the receiverID of the message. This is a common field for
both PrivateMsg and GroupMsg. This is useful when we want to send a message to
the frontend, but we don't know if it is a PrivateMsg or GroupMsg.
*/
func (p *PrivateMsg) GetReceiverID() int {
	return p.ReceiverID
}

func (g *GroupMsg) GetReceiverID() int {
	return g.ReceiverID
}

/*
GetMessage() returns the message string. This is a common field for both
PrivateMsg and GroupMsg. This is useful when we want to send a message to
the frontend, but we don't know if it is a PrivateMsg or GroupMsg.
*/
func (p *PrivateMsg) GetMessage() string {
	return p.Message
}

func (g *GroupMsg) GetMessage() string {
	return g.Message
}

/*
GetTimestamp() returns the timestamp of the message. This is a common field for
both PrivateMsg and GroupMsg. This is useful when we want to send a message to
the frontend, but we don't know if it is a PrivateMsg or GroupMsg.
*/
func (p *PrivateMsg) GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (g *GroupMsg) GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*
GetMsgType() returns a string that identifies the type of message, and is a
common field for both PrivateMsg and GroupMsg. This is useful when we want to
send a message to the frontend, but we don't know if it is a PrivateMsg or
GroupMsg.
*/
func (p *PrivateMsg) GetMsgType() string {
	return "PrivateMsg"
}

func (g *GroupMsg) GetMsgType() string {
	return "GroupMsg"
}

/*
WrapMsg() returns a struct that contains the fields that are common to both
PrivateMsg and GroupMsg. This is useful when we want to send a message to
the frontend, but we don't know if it is a PrivateMsg or GroupMsg.
*/
func (p *PrivateMsg) WrapMsg() struct {
	Type    string `json:"type"`
	Payload struct {
		SenderID   int    `json:"senderID"`
		ReceiverID int    `json:"receiverID"`
		Message    string `json:"message"`
		Timestamp  string `json:"timeStamp"`
	} `json:"payload"`
} {
	return struct {
		Type    string `json:"type"`
		Payload struct {
			SenderID   int    `json:"senderID"`
			ReceiverID int    `json:"receiverID"`
			Message    string `json:"message"`
			Timestamp  string `json:"timeStamp"`
		} `json:"payload"`
	}{
		Type: "PrivateMsg",
		Payload: struct {
			SenderID   int    `json:"senderID"`
			ReceiverID int    `json:"receiverID"`
			Message    string `json:"message"`
			Timestamp  string `json:"timeStamp"`
		}{
			SenderID:   p.SenderID,
			ReceiverID: p.ReceiverID,
			Message:    p.Message,
			Timestamp:  p.GetTimestamp(),
		},
	}
}

func (g *GroupMsg) WrapMsg() struct {
	Type    string `json:"type"`
	Payload struct {
		SenderID   int    `json:"senderID"`
		ReceiverID int    `json:"receiverID"`
		Message    string `json:"message"`
		Timestamp  string `json:"timeStamp"`
	} `json:"payload"`
} {
	return struct {
		Type    string `json:"type"`
		Payload struct {
			SenderID   int    `json:"senderID"`
			ReceiverID int    `json:"receiverID"`
			Message    string `json:"message"`
			Timestamp  string `json:"timeStamp"`
		} `json:"payload"`
	}{
		Type: "GroupMsg",
		Payload: struct {
			SenderID   int    `json:"senderID"`
			ReceiverID int    `json:"receiverID"`
			Message    string `json:"message"`
			Timestamp  string `json:"timeStamp"`
		}{
			SenderID:   g.SenderID,
			ReceiverID: g.ReceiverID,
			Message:    g.Message,
			Timestamp:  g.GetTimestamp(),
		},
	}
}
