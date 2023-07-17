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
InsertEvent inserts an given event into the database and returns an error if it fails
*/
func InsertEvent(event db.Event) error {
	_, err := db.InsertData("events", event.CreatorId, event.GroupId, event.Title, event.Content, time.Now(), event.Date)
	if err != nil {
		return errors.New("Error inserting event" + err.Error())
	}
	return nil

}

/*
TODO : CreateEvent commenttts
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
ReadGroupEvents reads the events of a group from the database and creates a slice of GroupEvent structs,
which contain the event, the profile of the creator of the event, the participants of the event, and session id of the user,
using the ReadEventOptions function and the fillSmallProfile function,
and returns the slice and an error if it fails
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
			// TODO i don't think this is the best way to do this
			for _, u := range users["going"] {
				if u.UserId == userId {
					result[i].Status = "going"
					break
				}
			}
			if result[i].Status == "" {
				for _, u := range users["not_going"] {
					if u.UserId == userId {
						result[i].Status = "not_going"
						break
					}
				}
			}
		} else {
			return nil, fmt.Errorf("invalid event type at index %d", i)
		}
	}
	return result, nil
}

/*
TODO GetGroupEvents commenttts
*/
func GetGroupEvents(payload json.RawMessage) (Response, error) {
	var eventRequest Request

	err := json.Unmarshal(payload, &eventRequest)
	if err != nil {
		return Response{}, errors.New("Error unmarshalling event" + err.Error())
	}
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
