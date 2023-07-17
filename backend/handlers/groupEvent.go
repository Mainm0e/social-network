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
InsertEvent inserts an given event into events table in the database
and add a notification for each group member except the creator of the event,
it returns an error if any occurs otherwise it returns nil.
*/
func InsertEvent(event db.Event) error {
	_, err := db.InsertData("events", event.CreatorId, event.GroupId, event.Title, event.Content, time.Now(), event.Date)
	if err != nil {
		return errors.New("Error inserting event" + err.Error())
	}
	// insert group invitation in notifications table
	// get the group members
	members, err := db.FetchData("group_member", "groupId = ?", event.GroupId)
	if err != nil {
		return errors.New("Error fetching group members" + err.Error())
	}
	// insert notification for each member
	for _, m := range members {
		member := m.(db.GroupMember)
		if member.UserId != event.CreatorId {
			_, err = db.InsertData("notifications", member.UserId, event.CreatorId, event.GroupId, "new_event", time.Now())
			if err != nil {
				return errors.New("Error inserting group invitation" + err.Error())
			}
		}
	}
	return nil

}

/*
CreateEvent is a function that processes a creating group event request,
it gets the group event data from the payload, use the InsertEvent function to insert the event into the database,
and returns a response with a descriptive message and an event with the session id of the user. or an error if any occurred.
*/
func CreateEvent(payload json.RawMessage) (Response, error) {
	var event GroupEvent
	var response Response
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return Response{}, errors.New("Error unmarshalling event" + err.Error())
	}
	err = InsertEvent(event.Event)
	if err != nil {
		return Response{}, errors.New("Error inserting event" + err.Error())
	}
	payload, err = json.Marshal(map[string]string{"sessionId": event.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	eventEvent := events.Event{
		Type:    "createEvent",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", eventEvent, http.StatusOK}, nil
}

/*
ReadEventOptions reads the participants of an event from the database and returns a map of the participants with key "going" or "not_going",
base on their option (going/not-going) and an error if it fails
*/
func ReadEventOptions(eventId int) (map[string][]SmallProfile, error) {
	options, err := db.FetchData("event_member", "eventId = ?", eventId)
	if err != nil {
		return nil, errors.New("Error fetching event options" + err.Error())
	}
	if len(options) == 0 {
		return nil, nil
	}
	result := make(map[string][]SmallProfile, len(options))
	for i, o := range options {
		if option, ok := o.(db.EventMember); ok {
			if option.Option == "going" {
				user, err := fillSmallProfile(option.MemberId)
				result["going"] = append(result["going"], user)
				if err != nil {
					return nil, errors.New("Error fetching event option" + err.Error())
				}
			} else if option.Option == "not_going" {
				user, err := fillSmallProfile(option.MemberId)
				result["not_going"] = append(result["not_going"], user)
				if err != nil {
					return nil, errors.New("Error fetching event option" + err.Error())
				}
			}
		} else {
			return nil, fmt.Errorf("invalid event option type at index %d", i)
		}
	}
	return result, nil
}

/*
eventUserStatus returns the option chosen by a user for an event, based on the map of the participants of the event and the user's id.
it used in ReadGroupEvents function and returns a string of the option chosen by the user("going"/"not_going")or an empty string if the user didn't choose an option.
*/
func eventUserStatus(users map[string][]SmallProfile, userId int) string {
	for _, u := range users["going"] {
		if u.UserId == userId {
			return "going"
		}
	}
	for _, u := range users["not_going"] {
		if u.UserId == userId {
			return "not_going"
		}
	}
	return ""
}

/*
ReadGroupEvents reads the events of a group from the database and creates a slice of GroupEvent structs,
which contain the event, the profile of the creator of the event, the participants of the event, session id of the user,
and the option chosen by the user for the event, using the ReadEventOptions function and the fillSmallProfile function, eventUserStatus function,
and returns the slice and an error if it fails otherwise it returns nil.
*/
func ReadGroupEvents(groupId int, userId int) ([]GroupEvent, error) {
	events, err := db.FetchData("events", "groupId = ?", groupId)
	if err != nil {
		return nil, errors.New("Error fetching group events" + err.Error())
	}

	if len(events) == 0 {
		return nil, errors.New("no events found")
	}

	result := make([]GroupEvent, len(events))
	for i, e := range events {
		if event, ok := e.(db.Event); ok {
			result[i].Event = event
			creatorProfile, err := fillSmallProfile(event.CreatorId)
			if err != nil {
				return nil, errors.New("Error fetching event creator profile" + err.Error())
			}

			result[i].CreatorProfile = creatorProfile
			users, err := ReadEventOptions(event.EventId)
			if err != nil {
				return nil, errors.New("Error fetching event options" + err.Error())
			}

			result[i].Participants = users
			// check this user's option
			result[i].Status = eventUserStatus(users, userId)
		} else {
			return nil, fmt.Errorf("invalid event type at index %d", i)
		}
	}
	return result, nil
}

/*
GetGroupEvents is a function that processes sending group events to the frontend,
it gets a payload containing the group id and the user id, use the ReadGroupEvents function to read the events of the group from the database,
and returns a response with a descriptive message and an event with the events payload. or an error if any occurred.
*/
func GetGroupEvents(payload json.RawMessage) (Response, error) {
	var eventRequest Request

	err := json.Unmarshal(payload, &eventRequest)
	if err != nil {
		return Response{}, errors.New("Error unmarshalling event" + err.Error())
	}
	// read group events from database
	groupEvents, err := ReadGroupEvents(eventRequest.GroupId, eventRequest.SenderId)
	if err != nil {
		return Response{}, errors.New("Error reading group events" + err.Error())
	}
	payload, err = json.Marshal(map[string][]GroupEvent{"events": groupEvents})
	if err != nil {
		return Response{}, errors.New("Error marshalling group events" + err.Error())
	}
	eventEvent := events.Event{
		Type:    "getGroupEvents",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", eventEvent, http.StatusOK}, nil

}

/*
InsertEventOption inserts the going/not-going of a group member to an event, into the database and returns an error if it fails
*/
func InsertEventOption(eventId int, memberId int, option string) error {
	// try to update the option if it already exists
	err := db.UpdateData("event_member", option, eventId, memberId)
	if err == nil {
		return nil
	}
	_, err = db.InsertData("event_member", eventId, memberId, option)
	if err != nil {
		return errors.New("Error inserting event option" + err.Error())
	}
	return nil
}

/*
ParticipateInEvent used for getting the going/not-going of a group member to an event from frontend and inserting it into the database,
using the InsertEventOption function, and returns an error if it fails
*/
func ParticipateInEvent(payload json.RawMessage) (Response, error) {
	var participateEvent db.EventMember
	err := json.Unmarshal(payload, &participateEvent)
	if err != nil {
		return Response{}, errors.New("Error unmarshalling event" + err.Error())
	}
	err = InsertEventOption(participateEvent.EventId, participateEvent.MemberId, participateEvent.Option)
	if err != nil {
		return Response{}, errors.New("Error inserting event option" + err.Error())
	}
	payload, err = json.Marshal(map[string]string{"sessionId": participateEvent.SessionId})
	if err != nil {
		return Response{}, errors.New("Error marshalling event option" + err.Error())
	}
	eventEvent := events.Event{
		Type:    "participateInEvent",
		Payload: payload,
	}
	return Response{"users retrieved successfully!", eventEvent, http.StatusOK}, nil
}
