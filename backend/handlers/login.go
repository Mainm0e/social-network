package handlers

import (
	"backend/events"
	"backend/server/sessions"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

/*
login is a function defined on the LoginData struct. It checks if user exists in the database base on the provided email,
and if the password matches the one stored in the database. If both are true,
it returns the user id and nil otherwise it returns 0 and an error.
*/
func (lg *LoginData) login() (int, error) {
	// Fetch user data from the database based on the provided email.
	user, err := fetchUser("email", lg.Email)
	if err != nil {
		return 0, errors.New("user not found")
	}
	// Compare the provided password with the password stored in the database.
	if user.Password == lg.Password {
		return user.UserId, nil
	} else {
		return 0, errors.New("password incorrect")
	}
}

/*
LoginPage is a function receiving a json file, payload, containing the login information.
It checks if the login information is valid, using the login() function. if it is,
it calls the sessions.Login() function to create a session for the user.
It then returns a Response struct containing the appropriate message and status code and an error if any occurred.
*/
func LoginPage(payload json.RawMessage) (Response, error) {
	var loginData LoginData
	err := json.Unmarshal(payload, &loginData)
	if err != nil {
		// handle the error
		log.Println("Error unmarshaling JSON to LoginData:", err)
	}
	log.Println("Login try by:", loginData)
	id, err := loginData.login()
	if err != nil {
		response := Response{err.Error(), events.Event{}, http.StatusBadRequest}
		log.Println("Error logging in:", err)
		return response, err
	} else {
		response := Response{"Login approved", events.Event{}, 200}
		sessionId, err := sessions.Login(loginData.Email, false)
		if err != nil {
			log.Println("Error logging in:", err)
			response := Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
		loginResponse := LoginResponse{sessionId, id}
		payload, err = json.Marshal(loginResponse)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, errors.New("Error marshaling profile to JSON: " + err.Error())
		}
		response.Event = events.Event{Type: "Login approved", Payload: payload}
		log.Println("Login approved", loginData)
		return response, nil
	}
}
