package handlers

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterPage(payload json.RawMessage) (Response, error) {
	var registerData RegisterData
	var response Response
	jsonData, err := json.Marshal(payload)
	if err != nil {
		// handle the error
		log.Println("Error marshaling payload to JSON:", err)
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}
	err = json.Unmarshal(jsonData, &registerData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to RegisterData:", err)
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}
	ok, err := IsNotUser(registerData.Email)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}
	if !ok {
		log.Println("User already exists:", registerData)
		response = Response{false, "User already exists", Event{}, http.StatusBadRequest}
		return response, err
	}
	// Check if an avatar image is provided
	if registerData.Avatar != "" {
		// Process the image and save it to the local storage
		url := "./images/avatars/" + registerData.Email
		err = utils.ProcessAvatarImage(registerData.Avatar, url)
		if err != nil {
			log.Println("Error processing avatar image:", err)
			response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
			return response, err
		}
		registerData.Avatar = url
	} else {
		registerData.Avatar = ""
	}
	err = registerData.register()
	if err != nil {
		log.Println("Error registering:", err)
		response = Response{false, "Registration denied", Event{}, 200}
		return response, err
	}
	response = Response{true, "Registration approved", Event{}, 200}
	log.Println("Registration approved:", registerData)
	return response, nil
}
