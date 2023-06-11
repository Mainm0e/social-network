package handlers

import "database/sql"

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterData struct {
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
