package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ProfilePage(payload json.RawMessage) (Response, error) {
	userId := 1 // TODO: get user id from session
	var response Response
	type UserId struct {
		UserID int `json:"user_id"`
	}
	var user UserId
	err := json.Unmarshal(payload, &user)
	if err != nil {
		// handle the error
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}

	var profile Profile
	profile, err = FillProfile(userId, user.UserID)
	if err != nil {
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, err
	}

	// Convert the Profile struct to JSON byte array
	payload, err = json.Marshal(profile)
	if err != nil {
		response = Response{false, err.Error(), Event{}, http.StatusBadRequest}
		return response, errors.New("Error marshaling profile to JSON: " + err.Error())
	}
	event := Event{
		Event_type: "profile",
		Payload:    payload,
	}
	response = Response{true, "profile data", event, http.StatusOK} // TODO: change message
	return response, nil

}
