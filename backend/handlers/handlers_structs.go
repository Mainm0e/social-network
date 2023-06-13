package handlers

import (
	"database/sql"
	"encoding/json"
)

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
type Response struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
type NullableString struct {
	sql.NullString
}

func (s *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.Valid = false
		return nil
	}
	s.Valid = true
	return json.Unmarshal(data, &s.String)
}
