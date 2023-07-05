package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

func InsertGroup(group Group) error {
	_, err := db.InsertData("groups", group.CreatorProfile.UserId, group.Title, group.Description, time.Now())
	if err != nil {
		return errors.New("Error inserting group" + err.Error())
	}
	return nil
}
func InsertGroupMember(groupId int, userId int) error {
	_, err := db.InsertData("group_member", groupId, userId)
	if err != nil {
		return errors.New("Error inserting group member" + err.Error())
	}
	return nil
}

/*
GetAllGroupMemberIDs returns all the user ids of the members of a group.
It takes a groupID integer as an argument and returns a slice of integers
and an error value, which is non-nil if any of the database operations
failed.
*/
func GetAllGroupMemberIDs(groupId int) ([]int, error) {
	var userIds []int

	// Fetch all group members from the database
	dbGroupMembers, err := db.FetchData("group_member", "groupId", groupId)
	if err != nil {
		return userIds, errors.New("Error fetching group members" + err.Error())
	}

	// If no group members were found, return an error
	if len(dbGroupMembers) == 0 {
		return userIds, errors.New("no group member found")
	}

	// Iterate over all group members and append their user ids to the slice
	for _, dbGroupMember := range dbGroupMembers {
		dbGroupMember := dbGroupMember.(db.GroupMember)
		userIds = append(userIds, dbGroupMember.UserId)
	}

	return userIds, nil
}

func ReadGroup(groupId int) (Group, error) {
	dbGroups, err := db.FetchData("groups", "groupId", groupId)
	if err != nil {
		return Group{}, errors.New("Error fetching group" + err.Error())
	}
	if len(dbGroups) == 0 {
		return Group{}, errors.New("group not found")
	}
	dbGroup := dbGroups[0].(db.Group)
	creator, err := fillSmallProfile(dbGroup.CreatorId)
	if err != nil {
		return Group{}, errors.New("Error fetching group creator" + err.Error())
	}
	group := Group{
		GroupId:        dbGroup.GroupId,
		CreatorProfile: creator,
		Title:          dbGroup.Title,
		Description:    dbGroup.Description,
		Date:           dbGroup.CreationTime,
	}
	return group, nil
}
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

func ExploreGroups(payload json.RawMessage) (Response, error) {
	var response Response
	var explore Explore
	err := json.Unmarshal(payload, &explore)
	log.Println("User: ", explore)
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
