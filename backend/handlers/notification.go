package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/*
readNotifications reads all notifications belonging to a user from the database. it gets the user id and session id as parameters,and
returns a slice of Notification structs, each containing a Notification data, the sender's public profile information(using fillSmallProfile function) , and user's SessionId.
and an error if any occurs otherwise it returns nil.
*/
func readNotifications(userId int, sessionId string) ([]Notification, error) {
	notifications, err := db.FetchData("notifications", "receiverId = ?", userId)
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

/*
RequestNotifications is a function that processes a getting Notification request by unmarshaling the payload,
validating the required fields, and calling readNotifications function to handle reading notifications from database.
It returns a response with a descriptive message and an event with the notifications payload. or an error if any occurred.
*/
func RequestNotifications(payload json.RawMessage) (Response, error) {
	var response Response
	var credential UserCredential
	err := json.Unmarshal(payload, &credential)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if credential.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if credential.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	notifications, err := readNotifications(credential.UserId, credential.SessionId)
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
