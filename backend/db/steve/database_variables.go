package database

import "database/sql"

// Formatting for terminal output
var Colour = Colours{
	Reset:    "\033[0m",
	Red:      "\033[31m",
	LightRed: "\033[1;31m",
	Orange:   "\033[0;33m",
	Yellow:   "\033[1;33m",
}

// Global database connection variable
var DB *sql.DB

/*
Global database insertion rules for ease of maintenance, and simplifying of the InsertData function
*/
var InsertRules = map[string]InsertRule{
	"users": {
		Query:      "INSERT INTO users(nickName, firstName, lastName, birthDate, email, password, creationTime) VALUES(?,?,?,?,?,?,?)",
		ExistTable: "users",
		ExistField: "nickName",
		ExistError: "nickname already exists",
	},
	"posts": {
		Query:          "INSERT INTO posts(userId, title, content, creationTime, status, groupId) VALUES(?,?,?,?,?,?)",
		NotExistTables: []string{"users", "groups"},
		NotExistFields: []string{"userId", "groupId"},
		NotExistErrors: []string{"user does not exist", "group does not exist"},
	},
	"comments": {
		Query:          "INSERT INTO comments(userId, postId, content, creationTime) VALUES(?,?,?,?)",
		NotExistTables: []string{"users", "posts"},
		NotExistFields: []string{"userId", "postId"},
		NotExistErrors: []string{"user does not exist", "post does not exist"},
	},
	"groups": {
		Query:          "INSERT INTO groups(creatorId, title, description) VALUES(?,?,?)",
		NotExistTables: []string{"users"},
		NotExistFields: []string{"creatorId"},
		NotExistErrors: []string{"creator does not exist"},
	},
	"messages": {
		Query:          "INSERT INTO messages(senderId, receiverId, messageContent, sendTime, seen) VALUES(?,?,?,?,?)",
		NotExistTables: []string{"users", "users"},
		NotExistFields: []string{"senderId", "receiverId"},
		NotExistErrors: []string{"sender does not exist", "receiver does not exist"},
	},
	"follow": {
		Query:          "INSERT INTO follow(followerId, followeeId, status) VALUES(?,?,?)",
		NotExistTables: []string{"users", "users"},
		NotExistFields: []string{"followerId", "followeeId"},
		NotExistErrors: []string{"follower does not exist", "followee does not exist"},
	},
	"group_member": {
		Query:          "INSERT INTO group_member(userId, groupId, status) VALUES(?,?,?)",
		NotExistTables: []string{"users", "groups"},
		NotExistFields: []string{"userId", "groupId"},
		NotExistErrors: []string{"user does not exist", "group does not exist"},
	},
	"semiPrivate": {
		Query:          "INSERT INTO semiPrivate(postId, userId) VALUES(?,?)",
		NotExistTables: []string{"posts", "users"},
		NotExistFields: []string{"postId", "userId"},
		NotExistErrors: []string{"post does not exist", "user does not exist"},
	},
	"notifications": {
		Query:          "INSERT INTO notifications(receiverId, senderId, type, content, creationTime) VALUES(?,?,?,?,?)",
		NotExistTables: []string{"users", "users"},
		NotExistFields: []string{"receiverId", "senderId"},
		NotExistErrors: []string{"receiver does not exist", "sender does not exist"},
	},
	"events": {
		Query:          "INSERT INTO events(creatorId, receiverId, groupId, title, content, creationTime, option) VALUES(?,?,?,?,?,?,?)",
		NotExistTables: []string{"users", "users", "groups"},
		NotExistFields: []string{"creatorId", "receiverId", "groupId"},
		NotExistErrors: []string{"creator does not exist", "receiver does not exist", "group does not exist"},
	},
}

/*
Global database table keys for ease of maintenance, and simplifying of the DeleteData function.
*/
var TableKeys = map[string]string{
	"users":     "userId",
	"posts":     "postId",
	"comments":  "commentId",
	"topics":    "topicId",
	"reactions": "reactionId",
}

/*
Global update rules for ease of maintenance, and simplifying of the UpdateData function.
*/
var UpdateRules = map[string]string{
	"users":      "UPDATE users SET userName=?, email=?, pass=? WHERE userId=?",
	"posts":      "UPDATE posts SET userId=?, title=?, content=? WHERE postId=?",
	"comments":   "UPDATE comments SET userId=?, postId=?, content=? WHERE commentId=?",
	"topics":     "UPDATE topics SET topicName=? WHERE topicId=?",
	"reactions":  "UPDATE reactions SET reaction=?, postId=?, commentId=? WHERE reactionId=?",
	"PostTopics": "UPDATE PostTopics SET postId=?, topicId=? WHERE postTopicId=?",
}

/*
FetchRules is a map of strings to structs. The strings are the names of the tables in the database,
and the structs are the rules for fetching data from the database. These anonymous structs contains
2 fields:
  - SelectFields: A string that specifies which columns to select when fetching data from
    the corresponding table.
  - ScanFields: A function that takes a pointer to a sql.Rows object and returns an interface{}
    and an error.

The map fascilitates ease of maintenance, and simplifying of the FetchData function.

When you use this map to fetch data from a table, you use the ScanFields function to scan the rows
returned by your query into the appropriate struct. The returned struct can then be used in your
application to process the data. This setup provides a lot of flexibility. For each table, you can
specify what data you want to fetch (SelectFields) and how to scan that data (ScanFields). And
because ScanFields is a function, it can contain any logic you need to correctly scan your rows
into the appropriate struct.
*/
var FetchRules = map[string]struct {
	SelectFields string
	/*
		ScanFields is a function that takes in a pointer to a sql.Rows object, and returns an interface
		and an error. The interface is the struct that the data will be scanned into, and the error is
		the error returned from the rows.Scan() function.
	*/
	ScanFields func(rows *sql.Rows) (interface{}, error)
}{
	"users": {
		SelectFields: "userId, nickName, email, pass, creationTime, age, gender, firstName, lastName, lastLogin",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var user User
			err := rows.Scan(&user.UserId, &user.NickName, &user.Email, &user.Pass, &user.CreationTime, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.LastLogin)
			return user, err
		},
	},
	"posts": {
		SelectFields: "postId, userId, topicId, title, content, creationTime",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var post Post
			err := rows.Scan(&post.PostId, &post.UserId, &post.TopicId, &post.Title, &post.Content, &post.CreationTime)
			return post, err
		},
	},
	"comments": {
		SelectFields: "commentId, userId, postId, content, creationTime",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var comment Comment
			err := rows.Scan(&comment.CommentId, &comment.UserId, &comment.PostId, &comment.Content, &comment.CreationTime)
			return comment, err
		},
	},
	"reactions": {
		SelectFields: "reactionId, userId, postId, commentId, reaction",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var reaction Reaction
			err := rows.Scan(&reaction.ReactionId, &reaction.UserId, &reaction.PostId, &reaction.CommentId, &reaction.Reaction)
			return reaction, err
		},
	},
	"topics": {
		SelectFields: "topicId, topicName",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var topic Topic
			err := rows.Scan(&topic.TopicId, &topic.TopicName)
			return topic, err
		},
	},
	"messages": {
		SelectFields: "messageId, senderId, receiverId, content, creationTime",
		ScanFields: func(rows *sql.Rows) (interface{}, error) {
			var message Message
			err := rows.Scan(&message.MessageId, &message.SenderId, &message.ReceiverId, &message.Content, &message.CreationTime)
			return message, err
		},
	},
}
