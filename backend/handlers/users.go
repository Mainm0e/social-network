package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"net/http"
)

func ExploreUsers(payload json.RawMessage) (Response, error) {
	var response Response
	var credential UserCredential
	err := json.Unmarshal(payload, &credential)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if credential.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if credential.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get users from database
	users, err := ReadAllUsers(credential.UserId, credential.SessionId)
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

/*
This function fetches a list of non-member users for a given group. It takes the groupId, userId, and sessionId as input parameters.
The function reads all users, retrieves the group members, and identifies the users who are not members of the group.
It also checks if the non-member users have already been invited to the group by the selected user.
The function returns the list of non-member users as a []Profile and any encountered errors.
*/
func nonMemberUsers(groupId int, userId int, sessionId string) ([]Profile, error) {
	usersIds, err := findFollowers(userId)
	if err != nil {
		return []Profile{}, errors.New("Error fetching users: " + err.Error())
	}

	members, err := db.FetchData("group_member", "groupId = ?", groupId)
	if err != nil {
		return []Profile{}, errors.New("Error fetching members: " + err.Error())
	}

	memberIds := make(map[int]struct{})
	for _, member := range members {
		if member.(db.GroupMember).Status == "member" || member.(db.GroupMember).Status == "pending" {
			memberIds[member.(db.GroupMember).UserId] = struct{}{}
		}
	}
	var nonMembers []Profile
	for _, id := range usersIds {
		_, exists := memberIds[id]
		if !exists {
			// now let see if this user already invited selected user to this group or not !!!!
			invitations, err := db.FetchData("notifications", "receiverId = ? AND senderId = ? AND groupId = ?", id, userId, groupId)
			if err != nil {
				return nil, errors.New("Error fetching notification data" + err.Error())
			}
			if len(invitations) == 0 {
				user, err := FillProfile(userId, id, sessionId)
				if err != nil {
					return nil, errors.New("Error fetching user" + err.Error())
				}

				// append the user to the list of non-members
				nonMembers = append(nonMembers, user)
			}
		}

	}
	return nonMembers, nil
}

func GetNonMembers(payload json.RawMessage) (Response, error) {
	var response Response
	var nonMembers NonMembers
	err := json.Unmarshal(payload, &nonMembers)
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
	users, err := nonMemberUsers(nonMembers.GroupId, nonMembers.UserId, nonMembers.SessionId)
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
		Type:    "getNonMembers",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", event, http.StatusOK}, nil
}
