package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func insertFollowRequest(senderId int, receiverId int) error {
	var reqType string
	var status string
	relation, err := checkUserRelation(senderId, receiverId)
	if err != nil {
		return errors.New("Error checking users relation" + err.Error())
	}
	user, err := fetchUser("userId", receiverId)
	if err != nil {
		return errors.New("Error fetching user" + err.Error())
	}
	if user.Privacy == "public" {
		reqType = "following"
		status = "following"
	} else {
		reqType = "follow_request"
		status = "pending"
	}
	if relation == "follow" {
		if user.Privacy == "private" {
			_, err := db.InsertData("notifications", receiverId, senderId, 0, reqType, time.Now())
			if err != nil {
				return errors.New("Error inserting follow request" + err.Error())
			}
		}
		_, err = db.InsertData("follow", senderId, receiverId, status)
		if err != nil {
			return errors.New("Error inserting follow request in follow table" + err.Error())
		}
	} else if relation == "pending" {
		notifications, err := db.FetchData("notifications", "senderId", senderId)
		if err != nil {
			return errors.New("Error fetching notifications" + err.Error())
		}
		for _, n := range notifications {
			if notification, ok := n.(db.Notification); ok {
				if notification.ReceiverId == receiverId && notification.Type == "follow_request" {
					err := db.DeleteData("notifications", notification.NotificationId)
					if err != nil {
						return errors.New("Error deleting notification" + err.Error())
					}
					break
				}
			}
		}
		err = db.UpdateData("follow", "following", senderId, receiverId)
		if err != nil {
			return errors.New("Error updating follow request" + err.Error())
		}
	} else if relation == "following" {

		err = db.DeleteData("follow", senderId, receiverId)
		if err != nil {
			return errors.New("Error deleting follow request" + err.Error())
		}
	}

	return nil

}

/*
DeleteFollowRequest function delete the follow request from notification table and update the follow table base on user decision
if error occur then it return error
*/

func deleteFollowRequest(followId int, notifId int, response string) error {
	err := db.DeleteData("notifications", followId)
	if err != nil {
		return errors.New("Error deleting follow request" + err.Error())
	}
	if response == "accept" {
		err = db.UpdateData("follow", "following", followId)
		if err != nil {
			return errors.New("Error updating follow request" + err.Error())
		}
	} else if response == "reject" {
		err = db.DeleteData("follow", followId)
		if err != nil {
			return errors.New("Error deleting follow request" + err.Error())
		}
	}
	return nil
}

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
	err = insertFollowRequest(follow.UserId, follow.FollowId)
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
	err = deleteFollowRequest(follow.UserId, follow.FollowId, follow.Response)
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
