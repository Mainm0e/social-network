package handlers

import (
	"backend/db"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	db.Check("../db/database.db", "./backend/database/test.sql")

	var tests = []struct {
		Email, Password any
		want            bool
	}{
		//test correct input
		{"john.doe@example.com", "password123", true},
		//test incorrect input
		{"john.doe@example.com", "password1234", false},
		//test incorrect input
		{"john.doe@example123.com", "password123", false},
		//test incorrect format input
		{23234, "password123", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.Email, tt.Password)
		t.Run(testname, func(t *testing.T) {
			// marshal the email and password into json base on LoginData struct and some random data
			data, err := json.Marshal(tt)
			if err != nil {
				t.Errorf("marshaling got error: %v:", err)
			}
			bo, err := login(data)

			if bo != tt.want {
				t.Errorf("login got: %v, want: %v error:%v.", bo, tt.want, err)
			}
		})

	}

}
func TestRegister(t *testing.T) {
	db.Check("../db/database.db", "./backend/database/test.sql")
	var tests = []struct {
		NickName     string
		FirstName    string
		LastName     string
		BirthDate    string
		Email        string
		Password     string
		AboutMe      string
		Avatar       string
		Privacy      string
		CreationTime string
		want         bool
	}{
		//test correct input
		{"johnny", "John", "Doe", "1990-05-15", "new.john.doe@example.com", "password123", "About John", "", "private", "2023-05-31 10:00:00", true},
		//test incorrect input (email already exists)
		{"johnny", "John", "Doe", "1990-05-15", "john.doe@example.com", "password123", "About John", "", "public", "2023-05-31 10:00:00", false},
	}
	for _, tt := range tests {

		testName := fmt.Sprintf("email: %s", tt.Email)
		t.Run(testName, func(t *testing.T) {
			// marshal the email and password into json base on LoginData struct and some random data
			data, err := json.Marshal(tt)
			if err != nil {
				t.Errorf("marshaling got error: %v:", err)
			}
			bo, err := register(data)

			if bo != tt.want {
				t.Errorf("login got: %v, want: %v error:%v.", bo, tt.want, err)
			}
		})

	}
}
