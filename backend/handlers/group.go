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
	if err != nil {
		return errors.New("Error inserting group" + err.Error())
	}
	group.GroupId = int(id)
	return nil
}

/*
InsertGroupMember inserts a new group member into "group_member" table. return error if any occurred otherwise returns nil.
*/
func InsertGroupMember(groupId int, userId int) error {
	_, err := db.InsertData("group_member", groupId, userId)
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
	err = InsertGroupMember(group.GroupId, group.CreatorProfile.UserId)
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
ReadGroup is defined as a method on the Group struct type.
it takes groupId as an argument fill the group struct with the group data from the database.
returns error if any occurred otherwise returns nil.
*/
func (group *Group) ReadGroup(groupId int) error {
	dbGroups, err := db.FetchData("groups", "groupId", groupId)
	if err != nil {
		return errors.New("Error fetching group" + err.Error())
	}
	if len(dbGroups) == 0 {
		return errors.New("group not found")
	}
	dbGroup := dbGroups[0].(db.Group)
	creator, err := fillSmallProfile(dbGroup.CreatorId)
	if err != nil {
		return errors.New("Error fetching group creator" + err.Error())
	}
	group.GroupId = dbGroup.GroupId
	group.CreatorProfile = creator
	group.Title = dbGroup.Title
	group.Description = dbGroup.Description
	group.Date = dbGroup.CreationTime

	return nil
}

/*
ReadAllGroups is a function that returns all existing groups in the database as a slice of Group struct.
it used in exploreGroups handler. and uses fillSmallProfile function to fill the creator profile.
if any error occurred it returns an error with a descriptive message.
*/
func ReadAllGroups(sessionId string) ([]Group, error) {
	dbGroups, err := db.FetchData("groups", "")
	if err != nil {
		return []Group{}, errors.New("Error fetching groups" + err.Error())
	}
	if len(dbGroups) == 0 {
		return []Group{}, errors.New("no group found")
	}
	var groups []Group
	for _, dbGroup := range dbGroups {
		dbGroup := dbGroup.(db.Group)
		creator, err := fillSmallProfile(dbGroup.CreatorId)
		if err != nil {
			return []Group{}, errors.New("Error fetching group creator" + err.Error())
		}
		group := Group{
			SessionId:      sessionId,
			GroupId:        dbGroup.GroupId,
			Title:          dbGroup.Title,
			Description:    dbGroup.Description,
			Date:           dbGroup.CreationTime,
			CreatorProfile: creator,
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
	//get groups from database
	groups, err := ReadAllGroups(explore.SessionId)
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
