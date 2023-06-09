package handlers

import (
	"backend/db"
	"encoding/json"
	"errors"
	"fmt"
)

func login(data []byte) (bool, error) {
	var login LoginData
	err := json.Unmarshal(data, &login)
	if err != nil {
		return false, errors.New("Error unmarshalling data" + err.Error())
	}
	user, err := db.FetchData("users", "email", login.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}
	if len(user) == 0 {
		return false, errors.New("user not found")
	}
	if user[0].(db.User).Password == login.Password {
		return true, nil
	} else {
		return false, errors.New("password incorrect")
	}
}
func register(data []byte) (bool, error) {
	var register RegisterData
	err := json.Unmarshal(data, &register)
	fmt.Println(register)
	if err != nil {
		return false, errors.New("Error unmarshalling data" + err.Error())
	}
	user, err := db.FetchData("users", "email", register.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}
	if len(user) == 0 {
		_, err := db.InsertData("users", register.NickName, register.FirstName, register.LastName, register.BirthDate, register.Email, register.Password, register.Privacy, register.CreationTime)
		if err != nil {
			return false, errors.New("Error inserting user" + err.Error())
		}
		return true, nil
	} else {
		return false, errors.New("user with this email already exists")
	}
}
