package handlers

import (
	"backend/events"
	"backend/server/sessions"
	"encoding/json"
	"log"
	"net/http"
)

func LogoutPage(payload json.RawMessage) (Response, error) {
	var logoutData LoginResponse
	err := json.Unmarshal(payload, &logoutData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to LogoutData:", err)
	}
	log.Println("Logout try by:", logoutData)
	err = sessions.Logout(logoutData.SessionId)
	if err != nil {
		response := Response{err.Error(), events.Event{}, http.StatusBadRequest}
		log.Println("Error logging out:", err)
		return response, err
	} else {
		response := Response{"Logout approved", events.Event{}, 200}
		return response, nil
	}
}
