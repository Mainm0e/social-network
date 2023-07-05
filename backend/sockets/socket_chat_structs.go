package sockets

type PrivateMsg struct {
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Message    string `json:"message"`
	Timestamp  string `json:"timeStamp"`
}

type GroupMsg struct {
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Message    string `json:"message"`
	Timestamp  string `json:"timeStamp"`
}

type ChatHistoryRequest struct {
	ChatType string `json:"chatType"`
	ClientID int    `json:"clientID"`
	TargetID int    `json:"targetID"` // UserID or GroupID
}

type ChatHistory struct {
	ChatType    string   `json:"chatType"`
	ClientID    int      `json:"clientID"`
	TargetID    int      `json:"targetID"` // UserID or GroupID
	ChatHistory []string `json:"chatHistory"`
}

type IsTyping struct {
	ChatType string `json:"chatType"`
	ClientID int    `json:"clientID"` // UserID of the client that is typing
	TargetID int    `json:"targetID"` // UserID / GroupID to identify chat in which typing is happening
}

// ChatMsg is an interface that is implemented by both PrivateMsg and GroupMsg.
// This allows us to use the same functions for both types of messages.
type ChatMsg interface {
	GetSenderID() int
	GetReceiverID() int
	GetMessage() string
	GetTimestamp() string
	GetType() string
}

/*********** NOTE TO SELF: NEVER FORGET HOW COOL INTERFACES ARE! ************/

func (p *PrivateMsg) GetSender() int {
	return p.SenderID
}

func (p *PrivateMsg) GetReceiver() int {
	return p.ReceiverID
}

func (p *PrivateMsg) GetMessage() string {
	return p.Message
}

func (p *PrivateMsg) GetTimestamp() string {
	return p.Timestamp
}

func (g *GroupMsg) GetSender() int {
	return g.SenderID
}

func (g *GroupMsg) GetReceiver() int {
	return g.ReceiverID
}

func (g *GroupMsg) GetMessage() string {
	return g.Message
}

func (g *GroupMsg) GetTimestamp() string {
	return g.Timestamp
}

func (p *PrivateMsg) GetChatType() string {
	return "PrivateMsg"
}

func (g *GroupMsg) GetChatType() string {
	return "GroupMsg"
}
