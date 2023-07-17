package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"fmt"
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
		_, err := db.InsertData("notifications", receiverId, senderId, 0, reqType, reqType, time.Now())
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
		notifications, err := db.FetchData("notifications", "senderId = ?", senderId)
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
insertGroupRequest is a function to manage a click on the "join"/"pending"/"leave?" button in a group page. // todo: change the name of the button base on frontend naming
based on the current relation between the user and the group. If the user is already a member of the group,
the function deletes the user from the group. If the user is not a member of the group, the function inserts a request
in the "notifications" table and a row in the "group_members" table with the status "pending." If the user has a invitation to the group already, //TODO: Idk what to do with this case ,again, it should be base on frontend naming
Returns nil on success; otherwise, returns an error with a descriptive message.
*/
func insertGroupRequest(senderId int, groupId int) error {
	dbGroups, err := db.FetchData("groups", "groupId = ?", groupId)
	if err != nil {
		return errors.New("Error fetching groups" + err.Error())
	}
	// getting group creator id
	receiverId := dbGroups[0].(db.Group).CreatorId
	if receiverId == 0 {
		return errors.New("error fetching group creator")
	}
	//getting relation between sender and group ("member", "pending", "waiting", "join")
	relation, err := groupUserRelation(senderId, groupId)
	if err != nil {
		return errors.New("Error checking users relation" + err.Error())
	}
	// dealing with request based on the relation type
	switch relation {

	case "join":
		//inserting group request in notifications table
		_, err := db.InsertData("notifications", receiverId, senderId, groupId, dbGroups[0].(db.Group).Title, "group_request", time.Now())
		if err != nil {
			return errors.New("Error inserting group request in notifications: " + err.Error())
		}
		//inserting group request in group_members table
		_, err = db.InsertData("group_member", senderId, groupId, "pending")
		if err != nil {
			return errors.New("Error inserting group request in group_member:" + err.Error())
		}
	case "member":
		//delete group member from group_members table
		err = db.DeleteData("group_member", groupId, senderId)
		if err != nil {
			return errors.New("Error deleting group member" + err.Error())
		}
	case "pending":
		//delete group request from notifications table
		notifications, err := db.FetchData("notifications", "groupId = ? AND senderId = ? AND type = ?", groupId, senderId, "group_request")
		if err != nil {
			return errors.New("Error fetching notifications" + err.Error())
		}
		err = db.DeleteData("notifications", notifications[0].(db.Notification).NotificationId)
		if err != nil {
			return errors.New("Error deleting notification" + err.Error())
		}
		//delete user from group_members table
		err = db.DeleteData("group_member", groupId, senderId)
		if err != nil {
			return errors.New("Error deleting group member" + err.Error())
		}
	case "waiting":
		fmt.Println("waiting")
		//TODO: IDK what to do here yet
	}
	return nil
}

/*
insertGroupInvitation is a function to manage a click on the "invite" button in a group page.
If the user is already invited to the group by same sender, the function does nothing. If the user is not invited to the group,
the function inserts a request in the "notifications" table and a row in the "group_members" table with the status "waiting" for the receiver.
Returns nil on success; otherwise, returns an error with a descriptive message.
*/
func insertGroupInvitation(senderId int, groupId int, receiverId int) error {
	// checking if the receiver is already invited to the group by the specific sender
	notif, err := db.FetchData("notifications", "receiverId = ? AND senderId = ? AND groupId = ?", receiverId, senderId, groupId)
	if err != nil {
		return errors.New("Error fetching notifications" + err.Error())
	}
	// if the receiver is already invited to the group by the specific sender we return nil and do nothing
	if len(notif) != 0 {
		return nil
	}
	groups, err := db.FetchData("groups", "groupId = ?", groupId)
	if err != nil {
		return errors.New("Error fetching group title" + err.Error())
	}
	// insert group invitation in notifications table
	_, err = db.InsertData("notifications", receiverId, senderId, groups[0].(db.Group).Title, "group_invitation", time.Now())
	if err != nil {
		return errors.New("Error inserting group invitation" + err.Error())
	}
	// insert receiverId in group_members table with status "waiting"
	db.InsertData("group_member", receiverId, groupId, "waiting")
	if err != nil {
		return errors.New("Error inserting group invitation" + err.Error())
	}
	return nil
}

/*
FollowRequest is a function that processes a follow request by unmarshaling the payload,
validating the required fields, and calling insertFollowRequest function to handle followRequest.
It returns a response with success/failure status and an event containing sessionId.
*/
func FollowOrJoinRequest(payload json.RawMessage) (Response, error) {
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
	if follow.GroupId == 0 {
		err = insertFollowRequest(follow.SenderId, follow.ReceiverId, follow.NotifId)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else if follow.GroupId != 0 {

		if follow.ReceiverId != 0 {
			err = insertGroupInvitation(follow.SenderId, follow.GroupId, follow.ReceiverId)
			if err != nil {
				response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
				return response, err
			}
		} else {
			err = insertGroupRequest(follow.SenderId, follow.GroupId)
			if err != nil {
				response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
				return response, err
			}
		}
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

func deleteRequest(tableName string, userId int, receiverId int, notifId int, response string) error {
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
		err = db.UpdateData(tableName, status, userId, receiverId)
		if err != nil {
			return errors.New("Error updating request" + err.Error())
		}
	} else if response == "reject" {
		err = db.DeleteData(tableName, userId, receiverId)
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
func FollowOrJoinResponse(payload json.RawMessage) (Response, error) {
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
	if follow.GroupId == 0 {
		err = deleteRequest("follow", follow.SenderId, follow.ReceiverId, follow.NotifId, follow.Content)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else if follow.GroupId != 0 {

		if follow.ReceiverId != 0 {
			err = deleteRequest("group_member", follow.ReceiverId, follow.GroupId, follow.NotifId, follow.Content)
			if err != nil {
				response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
				return response, err
			}
		} else if follow.SenderId != 0 {
			err = deleteRequest("group_member", follow.SenderId, follow.GroupId, follow.NotifId, follow.Content)
			if err != nil {
				response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
				return response, err
			}
		}
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
