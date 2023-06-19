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

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func GetPost(payload json.RawMessage) (Response, error) {
	var response Response
	var request RequestPost
	err := json.Unmarshal(payload, &request)
	log.Println("User: ", request)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.PostId == 0 {
		response = Response{"postId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	// TODO: check if the userId is necessary to get from request
	if request.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}

	//get post from database
	post, err := ReadPost(request.PostId, request.UserId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	post.SessionId = request.SessionId
	payload, err = json.Marshal(post)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "requestPost",
		Payload: payload,
	}
	return Response{"post retrieved successfully!", event, http.StatusOK}, nil
}

func GetPosts(payload json.RawMessage) (Response, error) {
	var response Response
	request := map[string]any{}
	err := json.Unmarshal(payload, &request)
	log.Println("User: ", request)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request["sessionId"] == nil {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request["userId"] == nil {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get posts from database
	posts, err := ReadPost(request["postId"].(int), request["userId"].(int))
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(posts)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "requestPosts",
		Payload: payload,
	}
	return Response{"posts retrieved successfully!", event, http.StatusOK}, nil
}
