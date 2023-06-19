package handlers

import (
	"backend/events"
	"encoding/json"
	"log"
	"net/http"
)

func CreateComment(payload json.RawMessage) (Response, error) {
	var response Response
	var comment Comment
	err := json.Unmarshal(payload, &comment)
	log.Println("comment: ", comment)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}

	//insert new post into database
	err = InsertComment(comment)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//send back sessionId
	payload, err = json.Marshal(map[string]string{"sessionId": comment.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "createPost",
		Payload: payload,
	}

	response = Response{"comment created successfully!", event, http.StatusOK}
	return response, nil
}
