package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func ProfilePage(payload json.RawMessage) (Response, error) {
	userId := 1 // TODO: get user id from session
	var response Response

	var user ProfileRequest
	log.Println("Payload: ", string(payload))
	err := json.Unmarshal(payload, &user)
	log.Println("User: ", user)
	if err != nil {
		// handle the error
		response = Response{err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}

	var profile Profile
	profile, err = FillProfile(userId, user.UserId)
	if err != nil {
		response = Response{err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}

	// Convert the Profile struct to JSON byte array
	payload, err = json.Marshal(profile)
	if err != nil {
		response = Response{err.Error(), Event{}, http.StatusBadRequest}
		return response, errors.New("Error marshaling profile to JSON: " + err.Error())
	}
	event := Event{
		Event_type: "profile",
		Payload:    payload,
	}
	response = Response{"profile data", event, http.StatusOK} // TODO: change message
	return response, nil

}
