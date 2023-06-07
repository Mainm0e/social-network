package db

import "time"

type User struct {
	UserId       int       `json:"userId"`
	NickName     string    `json:"nickName"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	BirthDate    string    `json:"birthDate"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	AboutMe      string    `json:"aboutMe"`
	Avatar       string    `json:"avatar"`
	CreationTime time.Time `json:"creationTime"`
}

type Post struct {
	PostId       int       `json:"postId"`
	UserId       int       `json:"userId"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
	Status       string    `json:"status"`
	Image        string    `json:"image"`
	GroupId      int       `json:"groupId"`
}

type Comment struct {
	CommentId    int       `json:"commentId"`
	UserId       int       `json:"userId"`
	PostId       int       `json:"postId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}

type Reaction struct {
	ReactionId int    `json:"reactionId"`
	UserId     int    `json:"userId"`
	PostId     int    `json:"postId"`
	CommentId  int    `json:"commentId"`
	Reaction   string `json:"reaction"`
}

type Message struct {
	MessageId      int       `json:"messageId"`
	SenderId       int       `json:"senderId"`
	ReceiverId     int       `json:"receiverId"`
	MessageContent string    `json:"messageContent"`
	SendTime       time.Time `json:"sendTime"`
	Seen           int       `json:"seen"`
}

type Follow struct {
	FollowerId int    `json:"followerId"`
	FolloweeId int    `json:"followeeId"`
	Status     string `json:"status"`
}

type Group struct {
	GroupId     int    `json:"groupId"`
	CreatorId   int    `json:"creatorId"`
	Title       string `json:"title"`
	Description string `json:"description"`
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
	NotificationId int       `json:"notificationId"`
	ReceiverId     int       `json:"receiverId"`
	SenderId       int       `json:"senderId"`
	Type           string    `json:"type"`
	Content        string    `json:"content"`
	CreationTime   time.Time `json:"creationTime"`
}

type Event struct {
	EventId      int       `json:"eventId"`
	CreatorId    int       `json:"creatorId"`
	ReceiverId   int       `json:"receiverId"`
	GroupId      int       `json:"groupId"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
	Option       string    `json:"option"`
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
