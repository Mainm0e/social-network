package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func readNotifications(userId int, sessionId string) ([]Notification, error) {
	notifications, err := db.FetchData("notifications", "receiverId", userId)
	if err != nil {
		return []Notification{}, errors.New("Error fetching notifications" + err.Error())
	}
	result := make([]Notification, len(notifications))
	for i, n := range notifications {
		if notification, ok := n.(db.Notification); ok {
			profile, err := fillSmallProfile(notification.SenderId)
			if err != nil {
				return []Notification{}, errors.New("Error fetching profile" + err.Error())
			}
			result[i] = Notification{
				SessionId:    sessionId,
				Notification: notification,
				Profile:      profile,
			}
		} else {
			return nil, fmt.Errorf("invalid notification type at index %d", i)
		}
	}

	return result, nil
}
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
	notifications, err := readNotifications(request.UserId, request.SessionId)
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
