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
 */
func InsertEventOption(eventId int, memberId int, option string) error {
	_, err := db.InsertData("event_member", eventId, memberId, option)
	if err != nil {
		return errors.New("Error inserting event option" + err.Error())
	}
	return nil
}
func ReadEventOptions(eventId int) (map[string][]SmallProfile, error) {
	options, err := db.FetchData("event_member", "eventId = ?", eventId)
	if err != nil {
		return nil, errors.New("Error fetching event options" + err.Error())
	}
	if len(options) == 0 {
		return nil, errors.New("no options found")
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
func ReadGroupEvents(groupId int) ([]GroupEvent, error) {
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
			users, err := ReadEventOptions(event.EventId)
			if err != nil {
				return nil, errors.New("Error fetching event options" + err.Error())
			}
			result[i].Participants = users
		} else {
			return nil, fmt.Errorf("invalid event type at index %d", i)
		}
	}
	return result, nil
}
func CreateEvent(payload json.RawMessage) (Response, error) {
	var event GroupEvent
	var response Response
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return Response{}, errors.New("Error unmarshalling event" + err.Error())
	}
	fmt.Println("event: ", event)
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
TODO:
func addEvent(payload json.RawMessage) (Response, error)     {}
func SendEvents(payload json.RawMessage) (Response, error)   {}
func EventOptions(payload json.RawMessage) (Response, error) {}
*/
// TODO: not sure but maybe we need to add a function to send all users participating in an event
