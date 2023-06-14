package handlers

type Event struct {
	Event_type string         `json:"event_type"`
	Payload    map[string]any `json:"payload"`
}

var Events = map[string]func(map[string]any) (Response, error){
	"login":    LoginPage,
	"register": RegisterPage,
	//"createPost": CreatePost,
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterData struct {
	NickName  string `json:"nickName,omitempty"` // optional
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthdate"`
	Email     string `json:"email"`
	Password  string `json:"matchPassword"`
	AboutMe   string `json:"aboutme,omitempty"` // optional
	Avatar    string `json:"avatar,omitempty"`  // optional
}

// TODO: we could remove success and make message more general
type Response struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// add post struct coming from frontend
/* type Post struct {
	UserId  int    `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"` ------> this one is important if its semi-private we need to get those followers id too and should handle in frontend that if its semi-private then user have to select followers.
	followers []int `json:"followers"`---> this one related to status
	Image   string `json:"image"`
	GroupId int    `json:"groupId"` ---> if post is a group post
}
	comment struct coming from frontend
	type Comment struct {
		PostId   int    `json:"postId"`
		UserId   int    `json:"userId"`
		Content  string `json:"content"`
		ParentId int    `json:"parentId"`
	}

	send post struct to frontend
	type SendPost struct {
	PostId   int    `json:"postId"`
	Post    Post   `json:"post"`
	Comments []Comment `json:"comments"`
	}

*/
