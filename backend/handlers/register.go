package handlers

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterPage(payload map[string]any) (Response, error) {
	// Retrieve the uploaded image file
	var registerData RegisterData
	var response Response
	jsonData, err := json.Marshal(payload)
	if err != nil {
		// handle the error
		fmt.Println("Error marshaling payload to JSON:", err)

	}
	err = json.Unmarshal(jsonData, &registerData)
	if err != nil {
		// handle the error
		fmt.Println("Error unmarshaling JSON to LoginData:", err)

	}
	err = registerData.register()
	if err != nil {
		log.Println("Error registering:", err)
		response = Response{false, "Registration denied", 200}
		return response, err
	}
	// Check if an avatar image is provided
	if registerData.Avatar != "" {
		// Process the image and save it to the local storage
		url := "./images/avatars/" + registerData.Email
		err = utils.ProcessAvatarImage(registerData.Avatar, url)
		if err != nil {
			log.Println("Error processing avatar image:", err)
			response = Response{false, err.Error(), http.StatusBadRequest}
			return response, err
		}
	}
	response = Response{true, "Registration approved", 200}
	log.Println("Registration approved:", registerData)
	return response, nil
}
