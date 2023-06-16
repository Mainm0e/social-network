package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func LoginPage(payload json.RawMessage) (Response, error) {
	var loginData LoginData
	err := json.Unmarshal(payload, &loginData)
	if err != nil {
		// handle the error
		fmt.Println("Error unmarshaling JSON to LoginData:", err)

	}
	log.Println("Login try by:", loginData)
	_, err = loginData.login()
	if err != nil {
		response := Response{err.Error(), Event{}, http.StatusBadRequest}
		log.Println("Error logging in:", err)
		return response, err
	} else {
		response := Response{"Login approved", Event{}, 200}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
