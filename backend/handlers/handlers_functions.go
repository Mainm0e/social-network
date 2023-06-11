package handlers

import (
	"backend/db"
	"encoding/json"
	"errors"
)

/*
login is a function that attempts to log in a user based on the provided data.
It takes in a byte slice `data` containing the login information.
It returns a boolean value indicating whether the login was successful, and an error if any occurred.
*/
func login(data []byte) (bool, error) {
	// Unmarshal the data into a LoginData struct.
	var login LoginData
	err := json.Unmarshal(data, &login)
	if err != nil {
		return false, errors.New("Error unmarshalling data" + err.Error())
	}

	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", login.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email was found.
	if len(user) == 0 {
		return false, errors.New("user not found")
	}

	// Compare the provided password with the password stored in the database.
	if user[0].(db.User).Password == login.Password {
		return true, nil
	} else {
		return false, errors.New("password incorrect")
	}
}

/*
register is a function that attempts to register a new user based on the provided data.
It takes in a byte slice `data` containing the registration information.
It returns a boolean value indicating whether the registration was successful, and an error if any occurred.
*/
func register(data []byte) (bool, error) {
	// Unmarshal the data into a RegisterData struct.
	var register RegisterData
	err := json.Unmarshal(data, &register)
	if err != nil {
		return false, errors.New("Error unmarshalling data" + err.Error())
	}
	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", register.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email already exists.
	if len(user) == 0 {
		// Insert the new user data into the database.
		_, err := db.InsertData("users", register.NickName, register.FirstName, register.LastName, register.BirthDate, register.Email, register.Password, register.AboutMe, register.Avatar, register.Privacy, register.CreationTime)
		if err != nil {
			return false, errors.New("Error inserting user" + err.Error())
		}
		return true, nil
	} else {
		return false, errors.New("user with this email already exists")
	}
}
