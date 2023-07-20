package handlers

import (
	"backend/events"
	"backend/server/sessions"
	"encoding/json"
	"log"
	"net/http"
)

/*
Logout is a function receiving a json file, payload, containing userCredential (sessionId and userId).
it calls the sessions.Logout() function to delete a user session in backend.
It then returns a Response struct containing the appropriate message and status code and an error if any occurred.
*/
func Logout(payload json.RawMessage) (Response, error) {
	var logoutData UserCredential
	err := json.Unmarshal(payload, &logoutData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to LogoutData:", err)
	}
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
