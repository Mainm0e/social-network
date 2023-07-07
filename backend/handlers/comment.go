package handlers

import (
	"backend/db"
	"backend/events"
	"backend/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
InsertComment inserts a comment into the database base on the Comment struct
if comment has an image it will be processed and saved to the local storage then the url will be saved to the database
if error occur it return error.
*/
func InsertComment(comment Comment) error {
	id, err := db.InsertData("comments", comment.UserId, comment.PostId, comment.Content, "", time.Now())
	if err != nil {
		return errors.New("Error inserting comment " + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting comment ")
	}
	if comment.Image != "" {
		// Process the image and save it to the local storage
		str := strconv.Itoa(int(id))
		url := "./images/comments/" + str
		url, err := utils.ProcessImage(comment.Image, url)
		if err != nil {
			log.Println("Error processing comment image:", err)
			//response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return err
		}
		comment.Image = url
	} else {
		comment.Image = ""
	}
	err = db.UpdateData("comments", comment.Image, id)
	if err != nil {
		return errors.New("Error updating comment image" + err.Error())
	}
	return nil
}

func CreateComment(payload json.RawMessage) (Response, error) {
	var response Response
	var comment Comment
	err := json.Unmarshal(payload, &comment)
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
