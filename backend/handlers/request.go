package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

/*
insertFollowRequest is a function to manage follow requests between users,
based on the privacy setting of the receiver and the current relation between the users.
If the receiver's privacy is "public," the follow request is established immediately as "following"
For "private" accounts, the function keeps the request as "follow_request" pending approval.
Handles insertions and deletions in "notifications" and "follow" tables.
Returns nil on success; otherwise, returns an error with a descriptive message.
*/
func insertFollowRequest(senderId int, receiverId int, notifId int) error {
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
	// use reqType to insert into notifications table and status to insert into follow table
	//reqType can be "follow_request" or "following"
	//status can be "pending" or "following"
	if user.Privacy == "public" {
		reqType = "following"
		status = "following"
	} else {
		reqType = "follow_request"
		status = "pending"
	}
	// if sender didn't request to follow the receiver before then insert the request in notifications and follow tables.
	if relation == "follow" {
		_, err := db.InsertData("notifications", receiverId, senderId, 0, reqType, time.Now())
		if err != nil {
			return errors.New("Error inserting follow request" + err.Error())
		}
		_, err = db.InsertData("follow", senderId, receiverId, status)
		if err != nil {
			return errors.New("Error inserting follow request in follow table" + err.Error())
		}
		// if sender requested to follow the receiver before then they's trying to take it back by click on follow button.
		// so we delete that follow_requset notification.
	} else if relation == "pending" {
		notifications, err := db.FetchData("notifications", "senderId", senderId)
		if err != nil {
			return errors.New("Error fetching notifications" + err.Error())
		}
		for _, n := range notifications {
			if notification, ok := n.(db.Notification); ok {
				// finding the follow request notification
				if notification.ReceiverId == receiverId && notification.Type == "follow_request" {
					err := db.DeleteData("notifications", notification.NotificationId)
					if err != nil {
						return errors.New("Error deleting notification" + err.Error())
					}
					break
				}
			}
		}
	}
	// if sender is following user or they have a pending relation in follow table we delete it.
	if relation == "following" || relation == "pending" {
		err = db.DeleteData("follow", senderId, receiverId)
		if err != nil {
			return errors.New("Error deleting follow request" + err.Error())
		}
	}

	return nil

}

/*
FollowRequest is a function that processes a follow request by unmarshaling the payload,
validating the required fields, and calling insertFollowRequest function to handle followRequest.
It returns a response with success/failure status and an event containing sessionId.
*/
func FollowRequest(payload json.RawMessage) (Response, error) {
	var response Response
	var follow Request
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
	if follow.SenderId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.ReceiverId == 0 {
		response = Response{"followId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = insertFollowRequest(follow.SenderId, follow.ReceiverId, follow.NotifId)
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

/*
DeleteRequest function delete the follow/join-group request from notification table and update the follow/group_member table base on user decision
if error occur then it return error
*/

func deleteRequest(tableName string, senderId int, receiverId int, notifId int, response string) error {
	err := db.DeleteData("notifications", notifId)
	if err != nil {
		return errors.New("Error deleting request" + err.Error())
	}
	var status string
	if response == "accept" {
		if tableName == "follow" {
			status = "following"
		} else {
			status = "member"
		}
		err = db.UpdateData(tableName, status, senderId, receiverId)
		if err != nil {
			return errors.New("Error updating request" + err.Error())
		}
	} else if response == "reject" {
		err = db.DeleteData(tableName, senderId, receiverId)
		if err != nil {
			return errors.New("Error deleting  request" + err.Error())
		}
	}
	return nil
}

/*
FollowResponse is a function that processes a response to follow request/following notification by unmarshaling the payload,
validating the required fields, and calling deleteFollowRequest function to handle response and delete the notification.
It returns a response with success/failure status and an event containing sessionId.
*/
func FollowResponse(payload json.RawMessage) (Response, error) {
	var response Response
	var follow Request
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
	if follow.SenderId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if follow.ReceiverId == 0 {
		response = Response{"followId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = deleteRequest("follow", follow.SenderId, follow.ReceiverId, follow.NotifId, follow.Content)
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

/*
group invitation && group request to join
*/
func insertGroupInvitation(senderId int, groupId int, receiverId int, content string) error {
	_, err := db.InsertData("notifications", receiverId, senderId, groupId, "group_invitation", time.Now())
	if err != nil {
		return errors.New("Error inserting group invitation" + err.Error())
	}
	return nil
}
func SendInvitation(payload json.RawMessage) (Response, error) {
	var response Response
	var invitation Request
	err := json.Unmarshal(payload, &invitation)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if invitation.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = insertGroupInvitation(invitation.SenderId, invitation.GroupId, invitation.ReceiverId, invitation.Content)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(map[string]string{"sessionId": invitation.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "sendInvitation",
		Payload: payload,
	}
	response = Response{"invitation sent successfully!", event, http.StatusOK}
	return response, nil
}
