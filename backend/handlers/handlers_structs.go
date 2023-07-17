package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
)

var Events = map[string]func(json.RawMessage) (Response, error){
	"login":          LoginPage,
	"register":       RegisterPage,
	"logout":         Logout,
	"profile":        ProfilePage,
	"updatePrivacy":  UpdatePrivacy,
	"profileList":    ProfileList,
	"createPost":     CreatePost,
	"getPosts":       GetPosts, // TODO: Change spelling / syntax
	"createComment":  CreateComment,
	"exploreUsers":   ExploreUsers,
	"followRequest":  FollowOrJoinRequest,
	"followResponse": FollowOrJoinResponse,
	"requestNotif":   RequestNotifications,
	"createGroup":    CreateGroup,
	"exploreGroups":  ExploreGroups,
	"getNonMembers":  GetNonMembers,
	"createEvent":    CreateEvent,
	"getGroupEvents": GetGroupEvents,
	"participate":    ParticipateInEvent,
}

type Response struct {
	Message    string       `json:"message"`
	Event      events.Event `json:"event"`
	StatusCode int          `json:"statusCode"`
}
type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserCredential struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
}

type ProfileListRequest struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	Request   string `json:"request"`
}

type ProfileRequest struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	ProfileId int    `json:"profileId"`
}
type RequestPost struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	PostId    int    `json:"postId"`
}
type RegisterData struct {
	NickName  string `json:"nickName,omitempty"` // optional
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthdate"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AboutMe   string `json:"aboutme,omitempty"` // optional
	Avatar    string `json:"avatar,omitempty"`  // optional
}

type SmallProfile struct {
	UserId    int     `json:"userId"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Avatar    *string `json:"avatar"`
}

type Profile struct {
	SessionId    string         `json:"sessionId"`
	UserId       int            `json:"userId"`
	NickName     string         `json:"nickName"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Avatar       *string        `json:"avatar"` //
	Relation     string         `json:"relation"`
	Status       string         `json:"privacy"`
	FollowerNum  int            `json:"followerNum"`
	FollowingNum int            `json:"followingNum"`
	PrivateData  PrivateProfile `json:"privateProfile"`
}

type PrivateProfile struct {
	BirthDate string `json:"birthdate"`
	Email     string `json:"email"`
	AboutMe   string `json:"aboutme"`
	Followers []int  `json:"followers"`
	Following []int  `json:"following"`
}

type Comment struct {
	SessionId      string       `json:"sessionId"`
	CommentId      int          `json:"commentId"`
	PostId         int          `json:"postId"`
	UserId         int          `json:"userId"`
	CreatorProfile SmallProfile `json:"creatorProfile"`
	Content        string       `json:"content"`
	Image          string       `json:"image,omitempty"`
	Date           string       `json:"Date"`
}

type Post struct {
	SessionId      string       `json:"sessionId"`
	PostId         int          `json:"postId"`
	UserId         int          `json:"userId"`
	CreatorProfile SmallProfile `json:"creatorProfile"`
	Title          string       `json:"title"`
	Content        string       `json:"content"`
	Status         string       `json:"status"`    //---> if the post is "public", "private" or "semi-private"s.
	Followers      []int        `json:"followers"` //---> selected followers list, if its semi-private.
	Image          string       `json:"image,omitempty"`
	GroupId        int          `json:"groupId"` // ---> if post is a group post
	Comments       []Comment    `json:"comments"`
	Date           string       `json:"date"`
}

type ReqAllPosts struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	From      string `json:"from"` // from profile or from group
	ProfileId int    `json:"profileId"`
	GroupId   int    `json:"groupId"`
}

type Request struct {
	SessionId  string `json:"sessionId"`
	SenderId   int    `json:"senderId"`
	ReceiverId int    `json:"receiverId"`
	GroupId    int    `json:"groupId"`
	NotifId    int    `json:"notifId"`
	Content    string `json:"content"`
}
type Group struct {
	SessionId      string         `json:"sessionId"`
	CreatorProfile SmallProfile   `json:"creatorProfile"`
	GroupId        int            `json:"groupId"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Status         string         `json:"status"` // join, pending, waiting, member
	NoMembers      int            `json:"noMembers"`
	Members        []SmallProfile `json:"members"`
	Date           string         `json:"date"`
}

type GroupEvent struct {
	SessionId      string                    `json:"sessionId"`
	CreatorProfile SmallProfile              `json:"creatorProfile"`
	Event          db.Event                  `json:"event"`
	Status         string                    `json:"status"` // going,not_going
	Participants   map[string][]SmallProfile `json:"participants"`
}

type Notification struct {
	SessionId    string          `json:"sessionId"`
	Profile      SmallProfile    `json:"profile"`
	Notification db.Notification `json:"notifications"`
}
type NonMembers struct {
	SessionId  string         `json:"sessionId"`
	UserId     int            `json:"userId"`
	GroupId    int            `json:"groupId"`
	NonMembers []SmallProfile `json:"nonMembers"`
}
