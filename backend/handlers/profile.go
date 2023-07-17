package handlers

import (
	"backend/db"
	"backend/events"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

/*
UpdateFollowRequests is called when the user change their privacy from private to public.
it gets the userId and update the follow notifications of that user in the database from pending to following.
if error occur it return error else it returns nil.
*/
func UpdateFollowRequests(userId int) error {
	notifications, err := db.FetchData("notifications", "receiverId = ? AND type = ?", userId, "follow_request")
	if err != nil {
		return errors.New("Error fetching notifications" + err.Error())
	}
	for _, n := range notifications {
		notification := n.(db.Notification)
		err := db.UpdateData("notifications", "following", notification.NotificationId)
		if err != nil {
			return errors.New("Error updating notification" + err.Error())
		}
		err = db.UpdateData("follow", "following", notification.SenderId, notification.ReceiverId)
		if err != nil {
			return errors.New("Error updating follow request" + err.Error())
		}

	}

	return nil
}

/*
UpdateProfile associated with the user clicked on the privacy button in frontend to change their privacy(public/private)
it gets the userId and update the privacy in the database base on the current privacy of user, if it was public it will be private and vice versa
in meanwhile it will update the follow requests to be accepted when user change their privacy from private to public.
in case of error it will return error otherwise nil.
*/
func updateProfile(userId int) error {
	var privacy string
	user, err := fetchUser("userId", userId)
	if err != nil {
		return errors.New("Error fetching user " + err.Error())
	}
	if user.Privacy == "public" {
		privacy = "private"
	} else {
		privacy = "public"
		// Update follow requests to be accepted when user change their privacy to public
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

/*
UpdatePrivacy associated with the user clicked on the privacy button in frontend to change their privacy(public/private)
it gets a payload which contains the credentials of the user who clicked on the privacy button.
it will call updateProfile function to update the privacy of the user.
in case of error it will return error otherwise it will return a response with a success message and an event with type updatePrivacy and payload of sessionId.
*/
func UpdatePrivacy(payload json.RawMessage) (Response, error) {
	var response Response
	var credential UserCredential
	err := json.Unmarshal(payload, &credential)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	err = updateProfile(credential.UserId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//send back sessionId
	payload, err = json.Marshal(map[string]string{"sessionId": credential.SessionId})
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

/*
ProfilePage associated with the user clicked on the profile of other users or them self to see the profile of the user.
it gets a payload which contains the credentials of the user who clicked on the profile and the profileId of the user who's profile is being viewed.
it will call FillProfile function to get the profile data of the user,base on the relation between current user and requested user.
in case of error it will return error otherwise it will return a response with a success message and an event with type profile and payload of profile data.
*/
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
	response = Response{"profile data", event, http.StatusOK}
	return response, nil

}

/*
ProfileList associated with the user clicked on the followers or followings in profile page to see the list of followers or followings.
it gets a payload which contains the credentials of the user who clicked on the followers or followings and the profileId of the user who's followers or followings is being viewed.
it will call SmallProfileList function to get the list of followers or followings of the user,base on the relation between current user and requested user.
in case of error it will return error otherwise it will return a response with a success message and an event with type profileList and payload of list of followers or followings.
*/
func ProfileList(payload json.RawMessage) (Response, error) {
	var user ProfileListRequest
	var response Response
	err := json.Unmarshal(payload, &user)
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
