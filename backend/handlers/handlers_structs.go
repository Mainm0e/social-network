package handlers

import (
	"backend/events"
	"encoding/json"
)

var Events = map[string]func(json.RawMessage) (Response, error){
	"login":         LoginPage,
	"register":      RegisterPage,
	"profile":       ProfilePage,
	"profileList":   ProfileList,
	"createPost":    CreatePost,
	"GetPost":       GetPost,
	"GetPosts":      GetPosts,
	"createComment": CreateComment,
	"exploreUsers":  ExploreUsers,
	"exploreGroups": ExploreGroups,
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
	FollowerNum  int            `json:"followerNum"`
	FollowingNum int            `json:"followingNum"`
	PrivateData  PrivateProfile `json:"privateProfile"`
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

type Group struct {
	SessionId      string       `json:"sessionId"`
	CreatorProfile SmallProfile `json:"creatorProfile"`
	GroupId        int          `json:"groupId"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Date           string       `json:"date"`
}
type Follow struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
	FollowId  int    `json:"followId"`
}
type Explore struct {
	SessionId string `json:"sessionId"`
	UserId    int    `json:"userId"`
}
