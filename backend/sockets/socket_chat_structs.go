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
	ClientUsername string `json:"clientUsername"` // Username of the clientthat is typing
	TargetName     string `json:"targetName"`     // Username / GroupTitle to identify chat in which typing is happening
}
