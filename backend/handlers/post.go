package handlers

import (
	"backend/events"
	"encoding/json"
	"log"
	"net/http"
)

func CreatePost(payload json.RawMessage) (Response, error) {
	var response Response
	var post Post
	err := json.Unmarshal(payload, &post)
	log.Println("User: ", post)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//insert new post into database
	err = InsertPost(post)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//send back sessionId
	payload, err = json.Marshal(map[string]string{"sessionId": post.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "createPost",
		Payload: payload,
	}

	response = Response{"post created successfully!", event, http.StatusOK}
	return response, nil
}
