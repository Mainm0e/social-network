package db

import (
	"database/sql"
)

type User struct {
	UserId       int             `json:"userId"`             // auto increment
	NickName     *sql.NullString `json:"nickName,omitempty"` // optional
	FirstName    string          `json:"firstName"`
	LastName     string          `json:"lastName"`
	BirthDate    string          `json:"birthDate"`
	Email        string          `json:"email"`
	Password     string          `json:"password"`
	AboutMe      *sql.NullString `json:"aboutMe,omitempty"` // optional
	Avatar       *string         `json:"avatar,omitempty"`  // optional
	Privacy      string          `json:"privacy"`           // default: public
	CreationTime string          `json:"creationTime"`
}
type Post struct {
	PostId       int    `json:"postId"`
	UserId       int    `json:"userId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	CreationTime string `json:"creationTime"`
	Status       string `json:"status"`
	Image        string `json:"image"`
	GroupId      int    `json:"groupId"`
}

type Comment struct {
	CommentId    int    `json:"commentId"`
	UserId       int    `json:"userId"`
	PostId       int    `json:"postId"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	CreationTime string `json:"creationTime"`
}

type Reaction struct {
	ReactionId int    `json:"reactionId"`
	UserId     int    `json:"userId"`
	PostId     int    `json:"postId"`
	CommentId  int    `json:"commentId"`
	Reaction   string `json:"reaction"`
}

type Message struct {
	MessageId      int    `json:"messageId"`
	SenderId       int    `json:"senderId"`
	ReceiverId     int    `json:"receiverId"`
	MessageContent string `json:"messageContent"`
	SendTime       string `json:"sendTime"`
	MsgType        string `json:"msgType"`
}

type Follow struct {
	FollowerId int    `json:"followerId"`
	FolloweeId int    `json:"followeeId"`
	Status     string `json:"status"`
}

type Group struct {
	GroupId      int    `json:"groupId"`
	CreatorId    int    `json:"creatorId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

type GroupMember struct {
	UserId  int    `json:"userId"`
	GroupId int    `json:"groupId"`
	Status  string `json:"status"`
}

type SemiPrivate struct {
	PostId int `json:"postId"`
	UserId int `json:"userId"`
}

type Notification struct {
	NotificationId int    `json:"notificationId"`
	ReceiverId     int    `json:"receiverId"`
	SenderId       int    `json:"senderId"`
	GroupId        int    `json:"groupId"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	CreationTime   string `json:"creationTime"`
}

type Event struct {
	EventId      int    `json:"eventId"`
	CreatorId    int    `json:"creatorId"`
	GroupId      int    `json:"groupId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Date         string `json:"date"`
	CreationTime string `json:"creationTime"`
}
type EventMember struct {
	SessionId string `json:"sessionId"`
	EventId   int    `json:"eventId"`
	MemberId  int    `json:"memberId"`
	Option    string `json:"option"` // going, not_going
}

type InsertRule struct {
	Query          string
	ExistTable     string
	ExistField     string
	ExistError     string
	NotExistTables []string
	NotExistFields []string
	NotExistErrors []string
}

type Colours struct {
	Reset      string // Resets terminal colour to default after 'text colouring'
	Red        string
	LightRed   string
	Green      string
	LightGreen string
	Blue       string
	LightBlue  string
	Orange     string
	Yellow     string
}
