package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

/*
Update Follow requests to be accepted when user change their privacy to public
*/
func UpdateFollowRequests(userId int) error {
	notifications, err := db.FetchData("notifications", "receiverId = ?", userId)
	if err != nil {
		return errors.New("Error fetching notifications" + err.Error())
	}
	for _, n := range notifications {
		if notification, ok := n.(db.Notification); ok {
			if notification.Type == "follow_request" {
				err := db.UpdateData("notifications", "following", notification.NotificationId)
				if err != nil {
					return errors.New("Error updating notification" + err.Error())
				}
				err = db.UpdateData("follow", "following", notification.SenderId, notification.ReceiverId)
				if err != nil {
					return errors.New("Error updating follow request" + err.Error())
				}
			}
		}
	}

	return nil
}

/*
UpdateProfile is a function that updates the privacy of a user with the specified email.
returns error if any occurred.
*/
func UpdateProfile(userId int, privacy string) error {
	user, err := fetchUser("userId", userId)
	if err != nil {
		return errors.New("Error fetching user " + err.Error())
	}
	fmt.Println("privacy:", privacy)
	if user.Privacy == "public" {
		privacy = "private"
	} else {
		privacy = "public"
		err = UpdateFollowRequests(userId)
		if err != nil {
			log.Println("Error updating follow requests", err.Error())
			return errors.New("Error updating follow requests " + err.Error())
		}
	}
	err = db.UpdateData("users", privacy, user.UserId)
	if err != nil {
		return errors.New("Error updating user " + err.Error())
	}
	return nil
}

func ProfilePage(payload json.RawMessage) (Response, error) {
	var response Response

	var user ProfileRequest
	err := json.Unmarshal(payload, &user)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}

	var profile Profile
	profile, err = FillProfile(user.UserId, user.ProfileId, user.SessionId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}

	// Convert the Profile struct to JSON byte array
	payload, err = json.Marshal(profile)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, errors.New("Error marshaling profile to JSON: " + err.Error())
	}
	event := events.Event{
		Type:    "profile",
		Payload: payload,
	}
	response = Response{"profile data", event, http.StatusOK} // TODO: change message
	return response, nil

}
func UpdatePrivacy(payload json.RawMessage) (Response, error) {
	var response Response
	var data PrivacyData
	err := json.Unmarshal(payload, &data)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = UpdateProfile(data.UserId, data.Privacy)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//send back sessionId
	payload, err = json.Marshal(map[string]string{"sessionId": data.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "updatePrivacy",
		Payload: payload,
	}
	response = Response{"privacy updated", event, http.StatusOK}
	return response, nil
}

func ProfileList(payload json.RawMessage) (Response, error) {
	var user ProfileListRequest
	var response Response
	log.Println("Payload: ", string(payload))
	err := json.Unmarshal(payload, &user)
	log.Println("User: ", user)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}

	list, err := SmallProfileList(user.UserId, user.Request)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	payload, err = json.Marshal(list)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, errors.New("Error marshaling profile to JSON: " + err.Error())
	}
	event := events.Event{
		Type:    "profileList",
		Payload: payload,
	}

	response = Response{"profile list", event, http.StatusOK}
	return response, nil
}
