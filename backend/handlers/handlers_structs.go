package handlers

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterData struct {
	NickName     string `json:"nickName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BirthDate    string `json:"birthDate"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AboutMe      string `json:"aboutMe"`
	Avatar       string `json:"avatar"`
	Privacy      string `json:"privacy"`
	CreationTime string `json:"creationTime"`
}
