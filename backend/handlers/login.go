package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func LoginPage(payload interface{}) (Response, error) {
	var loginData LoginData
	err := json.Unmarshal(payload.([]byte), &loginData)
	if err != nil {
		log.Println("Error decoding login data:", err)
		response := Response{false, err.Error(), http.StatusBadRequest}
		return response, err
	}
	log.Println("Login try by:", loginData)

	_, err = loginData.login()
	if err != nil {
		response := Response{false, err.Error(), http.StatusBadRequest}
		log.Println("Error logging in:", err)
		return response, err
	} else {
		response := Response{true, "Login approved", 200}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
