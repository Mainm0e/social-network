package handlers

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Retrieve the uploaded image file
		var registerData RegisterData
		var response Response

		err := json.NewDecoder(r.Body).Decode(&registerData)
		if err != nil {
			log.Println("Error decoding register data:", err)
			response = Response{false, err.Error(), http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		err = registerData.register()
		if err != nil {
			log.Println("Error registering:", err)
			response = Response{false, "Registration denied", 200}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Check if an avatar image is provided
		if registerData.Avatar != "" {
			// Process the image and save it to the local storage
			url := "./images/avatars/" + registerData.Email
			err = utils.ProcessAvatarImage(registerData.Avatar, url)
			if err != nil {
				log.Println("Error processing avatar image:", err)
				response = Response{false, err.Error(), http.StatusBadRequest}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		response = Response{true, "Registration approved", 200}
		log.Println("Registration approved:", registerData)
		json.NewEncoder(w).Encode(response)

	}
}
