package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
)

var Events = map[string]func(json.RawMessage) (Response, error){
	"login":    LoginPage,
	"register": RegisterPage,
	//TODO: "logout":         LogoutPage,
	"profile":        ProfilePage,
	"updatePrivacy":  UpdatePrivacy,
	"profileList":    ProfileList,
	"createPost":     CreatePost,
	"GetPosts":       GetPosts, // TODO: Change spelling / syntax
	"createComment":  CreateComment,
	"exploreUsers":   ExploreUsers,
	"followRequest":  FollowOrJoinRequest,
	"followResponse": FollowOrJoinResponse,
	"requestNotif":   RequestNotifications,
	"createGroup":    CreateGroup,
	"exploreGroups":  ExploreGroups,

	/*
		TODO: im not saying that we should have these functions but we need the functionality of these functions:
		"responseInvitation": AcceptInvitation,
		"requestToJoin":      RequestToJoin,
		"responseToJoin":     ResponseToJoin,
		"getGroupEvents":     GetGroupEvents,
		"responseEvent":      ResponseEvent,
	*/

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
type LoginResponse struct {
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

type PrivacyData struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	Privacy   string `json:"privacy"`
}

type PrivateProfile struct {
	BirthDate string `json:"birthdate"`
	Email     string `json:"email"`
	AboutMe   string `json:"aboutme"`
	Followers []int  `json:"followers"` // become array of uuid
	Following []int  `json:"following"` // become array of uuid
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
	Status         string       `json:"status"`    //------> this one is important if its semi-private we need to get those followers id too and should handle in frontend that if its semi-private then user have to select followers.
	Followers      []int        `json:"followers"` //---> this one related to status
	Image          string       `json:"image,omitempty"`
	GroupId        int          `json:"groupId"` // ---> if post is a group post
	Comments       []Comment    `json:"comments"`
	Date           string       `json:"date"`
}

type RequestPost struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	PostId    int    `json:"postId"`
}

type ReqAllPosts struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	From      string `json:"from"` // from profile or from group
	ProfileId int    `json:"profileId"`
	GroupId   int    `json:"groupId"`
}

type Explore struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
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
	SessionId    string                    `json:"sessionId"`
	Event        db.Event                  `json:"event"`
	Participants map[string][]SmallProfile `json:"participants"`
}

type Notification struct {
	SessionId    string          `json:"sessionId"`
	Profile      SmallProfile    `json:"profile"`
	Notification db.Notification `json:"notifications"`
}
