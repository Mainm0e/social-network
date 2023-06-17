package handlers

import (
	"backend/events"
	"backend/server/sessions"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func LoginPage(payload json.RawMessage) (Response, error) {
	var loginData LoginData
	err := json.Unmarshal(payload, &loginData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to LoginData:", err)
	}
	log.Println("Login try by:", loginData)
	_, err = loginData.login()
	if err != nil {
		response := Response{err.Error(), events.Event{}, http.StatusBadRequest}
		log.Println("Error logging in:", err)
		return response, err
	} else {
		response := Response{"Login approved", events.Event{}, 200}
		sessionId, err := sessions.Login(loginData.Email, false)
		if err != nil {
			log.Println("Error logging in:", err)
			response := Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
		payload, err = json.Marshal(map[string]string{"sessionId": sessionId})
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, errors.New("Error marshaling profile to JSON: " + err.Error())
		}
		response.Event = events.Event{"Login approved", payload}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
