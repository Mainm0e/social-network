package handlers

import (
	"backend/events"
	"encoding/json"
	"log"
	"net/http"
)

func ExploreGroups(payload json.RawMessage) (Response, error) {
	var response Response
	var explore Explore
	err := json.Unmarshal(payload, &explore)
	log.Println("User: ", explore)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if explore.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if explore.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get groups from database
	groups, err := ReadAllGroups(explore.SessionId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(groups)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "exploreGroups",
		Payload: payload,
	}
	return Response{"groups retrieved successfully!", event, http.StatusOK}, nil
}
