package handlers

import (
	"backend/db"
	"errors"
	"time"
)

/*
login is a function that attempts to log in a user based on the provided data.
It takes in a byte slice `data` containing the login information.
It returns a boolean value indicating whether the login was successful, and an error if any occurred.
*/
func (lg *LoginData) login() (bool, error) {

	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", lg.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email was found.
	if len(user) == 0 {
		return false, errors.New("user not found")
	}

	// Compare the provided password with the password stored in the database.
	if user[0].(db.User).Password == lg.Password {
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
func (regData *RegisterData) register() error {
	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", regData.Email)
	if err != nil {
		return errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email already exists.
	if len(user) == 0 {
		// Insert the new user data into the database.
		_, err := db.InsertData("users", regData.NickName, regData.FirstName, regData.LastName, regData.BirthDate, regData.Email, regData.Password, regData.AboutMe, regData.Avatar, "public", time.Now())
		if err != nil {
			return errors.New("Error inserting user" + err.Error())
		}
		return nil
	} else {
		return errors.New("user with this email already exists")
	}
}
