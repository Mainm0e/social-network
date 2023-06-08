package database

type User struct {
	UserId       int
	NickName     string
	Email        string
	Pass         string
	CreationTime string
	Age          int
	Gender       string
	FirstName    string
	LastName     string
	LastLogin    string
}

type Post struct {
	PostId       int
	UserId       int
	TopicId      int
	Title        string
	Content      string
	CreationTime string
}

type Topic struct {
	TopicName string
	TopicId   int
}

type Comment struct {
	CommentId    int
	UserId       int
	PostId       int
	Content      string
	CreationTime string
}

type Reaction struct {
	ReactionId int
	UserId     int
	PostId     int
	CommentId  int
	Reaction   string
}

type Message struct {
	MessageId    int
	SenderId     int
	ReceiverId   int
	Content      string
	CreationTime string
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
