package handlers

import (
	"backend/events"
	"encoding/json"
	"log"
	"net/http"
)

func ExploreUsers(payload json.RawMessage) (Response, error) {
	var response Response
	var explore Explore
	err := json.Unmarshal(payload, &explore)
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
	//get users from database
	users, err := ReadAllUsers(explore.UserId, explore.SessionId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(users)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "exploreUsers",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", event, http.StatusOK}, nil
}

func GetNonMembers(payload json.RawMessage) (Response, error) {
	var response Response
	var nonMembers NonMembers
	err := json.Unmarshal(payload, &nonMembers)
	log.Println("User: ", nonMembers)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if nonMembers.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if nonMembers.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get users from database
	users, err := NonMemberUsers(nonMembers.GroupId, nonMembers.UserId, nonMembers.SessionId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(users)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "exploreUsers",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", event, http.StatusOK}, nil
}
