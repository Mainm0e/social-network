package handlers

import (
	"backend/events"
	"encoding/json"
	"net/http"
)

func FollowRequest(payload json.RawMessage) (Response, error) {
	var response Response
	var follow Follow
	err := json.Unmarshal(payload, &follow)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.FollowId == 0 {
		response = Response{"followId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get users from database
	err = InsertFollowRequest(follow.UserId, follow.FollowId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(map[string]string{"sessionId": follow.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "followRequest",
		Payload: payload,
	}
	response = Response{"follow request sent successfully!", event, http.StatusOK}
	return response, nil
}
func FollowResponse(payload json.RawMessage) (Response, error) {
	var response Response
	var follow Follow
	err := json.Unmarshal(payload, &follow)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.FollowId == 0 {
		response = Response{"followId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = DeleteFollowRequest(follow.UserId, follow.FollowId, follow.Response)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(map[string]string{"sessionId": follow.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "followResponse",
		Payload: payload,
	}
	response = Response{"follow response sent successfully!", event, http.StatusOK}
	return response, nil
}
