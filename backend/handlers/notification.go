package handlers

import (
	"backend/events"
	"encoding/json"
	"net/http"
)

func RequestNotifications(payload json.RawMessage) (Response, error) {
	var response Response
	var request Explore
	err := json.Unmarshal(payload, &request)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	notifications, err := ReadNotifications(request.UserId, request.SessionId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(notifications)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "notifications",
		Payload: payload,
	}
	response = Response{"notifications", event, http.StatusOK}
	return response, nil
}
