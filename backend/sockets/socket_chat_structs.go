package sockets

type PrivateMsg struct {
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
	Message          string `json:"message"`
	Timestamp        string `json:"timeStamp"`
}

type GroupMsg struct {
	SenderUsername string `json:"senderUsername"`
	GroupTitle     string `json:"groupTitle"`
	Message        string `json:"message"`
	Timestamp      string `json:"timeStamp"`
}

type ChatHistoryRequest struct {
	ChatType       string `json:"chatType"`
	ClientUsername string `json:"clientUsername"`
	TargetName     string `json:"targetName"` // Username or GroupTitle
}

type ChatHistory struct {
	ChatType       string   `json:"chatType"`
	ClientUsername string   `json:"clientUsername"`
	TargetName     string   `json:"targetName"` // Username or GroupTitle
	ChatHistory    []string `json:"chatHistory"`
}

type IsTyping struct {
	ChatType       string `json:"chatType"`
	ClientUsername string `json:"clientUsername"` // Username of the client that is typing
	TargetName     string `json:"targetName"`     // Username / GroupTitle to identify chat in which typing is happening
}

// ChatMsg is an interface that is implemented by both PrivateMsg and GroupMsg.
// This allows us to use the same functions for both types of messages.
type ChatMsg interface {
	GetSender() string
	GetReceiver() string
	GetMessage() string
	GetTimestamp() string
	GetType() string
}

/*********** NOTE TO SELF: NEVER FORGET HOW COOL INTERFACES ARE! ************/

func (p *PrivateMsg) GetSender() string {
	return p.SenderUsername
}

func (p *PrivateMsg) GetReceiver() string {
	return p.ReceiverUsername
}

func (p *PrivateMsg) GetMessage() string {
	return p.Message
}

func (p *PrivateMsg) GetTimestamp() string {
	return p.Timestamp
}

func (g *GroupMsg) GetSender() string {
	return g.SenderUsername
}

func (g *GroupMsg) GetReceiver() string {
	return g.GroupTitle
}

func (g *GroupMsg) GetMessage() string {
	return g.Message
}

func (g *GroupMsg) GetTimestamp() string {
	return g.Timestamp
}

func (p *PrivateMsg) GetType() string {
	return "PrivateMsg"
}

func (g *GroupMsg) GetType() string {
	return "GroupMsg"
}
