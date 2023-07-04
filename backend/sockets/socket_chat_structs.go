package sockets

type PrivateMsgSend struct {
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
	Message          string `json:"message"`
	Timestamp        string `json:"timeStamp"`
}

type PrivateMsgBroadcast struct {
	SenderUsername   string `json:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername"`
	Message          string `json:"message"`
	Timestamp        string `json:"timeStamp"`
}

type GroupMsgSend struct {
	SenderUsername string `json:"senderUsername"`
	GroupTitle     string `json:"groupTitle"`
	Message        string `json:"message"`
	Timestamp      string `json:"timeStamp"`
}

type GroupMsgBroadcast struct {
	SenderUsername string `json:"senderUsername"`
	GroupTitle     string `json:"groupTitle"`
	Message        string `json:"message"`
	Timestamp      string `json:"timeStamp"`
}

type GetChatHistory struct {
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
	ClientUsername string `json:"clientUsername"`
	TargetName     string `json:"targetName"` // Username or GroupTitle
}
