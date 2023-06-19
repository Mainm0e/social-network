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
	id, err := loginData.login()
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
		loginResponse := LoginResponse{sessionId, id}
		payload, err = json.Marshal(loginResponse)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, errors.New("Error marshaling profile to JSON: " + err.Error())
		}
		response.Event = events.Event{Type: "Login approved", Payload: payload}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
