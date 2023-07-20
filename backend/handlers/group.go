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
InsertGroup is defined as a method on the Group struct type.
It inserts a new group into the database and fills the groupId field in the group struct with given groupId.
*/
func (group *Group) InsertGroup() error {
	id, err := db.InsertData("groups", group.CreatorProfile.UserId, group.Title, group.Description, time.Now())
	// print type of group.CreatorProfile.UserId
	if err != nil {
		return errors.New("Error inserting group" + err.Error())
	}
	group.GroupId = int(id)
	return nil
}

/*
InsertGroupMember inserts a new group member into "group_member" table. return error if any occurred otherwise returns nil.
*/
func InsertGroupMember(groupId int, userId int, status string) error {
	_, err := db.InsertData("group_member", userId, groupId, status)
	if err != nil {
		return errors.New("Error inserting group member" + err.Error())
	}
	return nil
}

/*
CreateGroup gets sessionId and creation group data as jsonified payload from frontend.
it calls CreateGroup function to insert group into database which after inserting group into database it fill groupId field in group struct.
then it calls InsertGroupMember function to insert creator of the group into group_member table.
it returns a response with a descriptive message and an event with the group payload. or an error if any occurred.
*/
func CreateGroup(payload json.RawMessage) (Response, error) {
	var response Response
	var group Group
	err := json.Unmarshal(payload, &group)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if group.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//insert group into database
	err = group.InsertGroup()
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//insert creator of the group into database
	err = InsertGroupMember(group.GroupId, group.CreatorProfile.UserId, "member")
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(map[string]string{"sessionId": group.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "createGroup",
		Payload: payload,
	}
	return Response{"group created successfully!", event, http.StatusOK}, nil
}

/*
groupUserRelation is a helper function that takes userId and groupId as an argument.
it returns the relation between the user and the group as a string.
base on "status" field in group_member table. if user does not have any relation with the group it returns "join".
(status could be one of these four values: "member", "pending", "waiting", "join")
if any error occurred it returns an error with a descriptive message.
*/
func groupUserRelation(userId, groupId int) (string, error) {
	groupMember, err := db.FetchData("group_member", "groupId = ?", groupId)
	if err != nil {
		return "", errors.New("Error fetching group member" + err.Error())
	}
	if len(groupMember) == 0 {
		return "", errors.New("group does not have any member")
	}
	for _, member := range groupMember {
		if member.(db.GroupMember).UserId == userId {
			return member.(db.GroupMember).Status, nil
		}
	}
	return "join", nil
}

/*
ReadGroup is defined as a method on the Group struct type. it takes db.Group struct and userId as an argument.
it fills the group struct with the group data from the database and other information like group members and group creator small profile.
it check the relation between the user and the group using groupUserRelation function.
returns error if any occurred otherwise returns nil.
*/
func (group *Group) ReadGroup(dbGroup db.Group, userId int) error {
	creator, err := fillSmallProfile(dbGroup.CreatorId)
	if err != nil {
		return errors.New("Error fetching group creator" + err.Error())
	}
	group.CreatorProfile = creator
	membersIds, err := GetAllGroupMemberIDs(group.GroupId)
	if err != nil {
		return errors.New("Error fetching group members" + err.Error())
	}

	membersProfiles := make([]SmallProfile, 0)
	for _, memberId := range membersIds {
		member, err := fillSmallProfile(memberId)
		if err != nil {
			return errors.New("Error fetching group member" + err.Error())
		}
		status, err := groupUserRelation(memberId, group.GroupId)
		if err != nil {
			return errors.New("Error fetching group member" + err.Error())
		}
		if status == "member" {
			membersProfiles = append(membersProfiles, member)
		}
	}
	group.Members = membersProfiles
	group.NoMembers = len(membersProfiles)

	status, err := groupUserRelation(userId, group.GroupId)
	if err != nil {
		return errors.New("Error fetching group member" + err.Error())
	}
	group.Status = status

	group.Title = dbGroup.Title
	group.Description = dbGroup.Description
	group.Date = dbGroup.CreationTime

	return nil
}

/*
ReadAllGroups is a function that returns all existing groups in the database as a slice of Group struct.
it used in exploreGroups handler. and uses ReadGroup function to fill a group data in the group struct,
if any error occurred it returns an error with a descriptive message.
*/
func ReadAllGroups(sessionId string, userId int) ([]Group, error) {
	dbGroups, err := db.FetchData("groups", "")
	if err != nil {
		return []Group{}, errors.New("Error fetching groups" + err.Error())
	}
	if len(dbGroups) == 0 {
		return []Group{}, errors.New("no group found")
	}
	var groups []Group
	for _, dbGroup := range dbGroups {
		var group Group
		group.GroupId = dbGroup.(db.Group).GroupId
		err := group.ReadGroup(dbGroup.(db.Group), userId)
		if err != nil {
			return []Group{}, errors.New("Error fetching group" + err.Error())
		}
		groups = append(groups, group)
	}
	return groups, nil
}

/*
ExploreGroups gets userId and sessionId as jsonified payload from frontend.
it calls ReadAllGroups function to get all groups from database.
it returns a response with a descriptive message and an event with the groups payload. or an error if any occurred.
*/
func ExploreGroups(payload json.RawMessage) (Response, error) {
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
	//get groups from database
	groups, err := ReadAllGroups(credential.SessionId, credential.UserId)
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
