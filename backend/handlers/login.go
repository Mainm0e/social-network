package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func LoginPage(payload json.RawMessage) (Response, error) {
	jsonData, err := json.Marshal(payload)
	log.Println("LoginPage payload:", string(jsonData))
	if err != nil {
		// handle the error
		fmt.Println("Error marshaling payload to JSON:", err)

	}
	var loginData LoginData
	err = json.Unmarshal(jsonData, &loginData)
	if err != nil {
		// handle the error
		fmt.Println("Error unmarshaling JSON to LoginData:", err)

	}
	log.Println("Login try by:", loginData)
	_, err = loginData.login()
	if err != nil {
		response := Response{false, err.Error(), Event{}, http.StatusBadRequest}
		log.Println("Error logging in:", err)
		return response, err
	} else {
		response := Response{true, "Login approved", Event{}, 200}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
