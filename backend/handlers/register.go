package handlers

import (
	"backend/db"
	"backend/events"
	"backend/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

/*
register is a function that attempts to register a new user based on the provided data.
It takes in a byte slice `data` containing the registration information.
It returns a boolean value indicating whether the registration was successful, and an error if any occurred.
*/
func (regData *RegisterData) register() error {
	_, err := db.InsertData("users", regData.Email, regData.FirstName, regData.LastName, regData.BirthDate, regData.NickName, regData.Password, regData.AboutMe, regData.Avatar, "private", time.Now())
	if err != nil {
		return errors.New("Error inserting user" + err.Error())
	}
	return nil
}

func RegisterPage(payload json.RawMessage) (Response, error) {
	var registerData RegisterData
	var response Response

	err := json.Unmarshal(payload, &registerData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to RegisterData:", err)
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	user, err := fetchUser("email", registerData.Email)
	if err != nil && err.Error() != "user not found" {
		log.Println("Error fetching user:", err)
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	} else if err == nil {
		log.Println("User already exists:", user)
		response = Response{"User already exists", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	// Check if an avatar image is provided
	if registerData.Avatar != "" {
		// Process the image and save it to the local storage
		url := "./images/avatars/" + registerData.Email
		url, err = utils.ProcessImage(registerData.Avatar, url)
		if err != nil {
			log.Println("Error processing avatar image:", err)
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
		registerData.Avatar = url
	} else {
		registerData.Avatar = ""
	}
	err = registerData.register()
	if err != nil {
		log.Println("Error registering:", err)
		response = Response{"Registration denied", events.Event{}, 200}
		return response, err
	}
	response = Response{"Registration approved", events.Event{}, 200}
	log.Println("Registration approved:", registerData)
	return response, nil
}
