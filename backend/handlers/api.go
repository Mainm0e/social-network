package handlers

import (
	"backend/events"
	"encoding/json"
	"log"
	"net/http"
)

/*
ApiHandler is the handler function for all endpoints. it's going to handle all the events that are sent through http requests. it's only function deal with http requests.
it decode the request body to an Event struct, then get the corresponding handler function for the event type and call it with the event payload as a parameter and get the response and error from it.
if there is an error, it will return a response with the error message and status code 400.
*/
func HTTPEventRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var event events.Event
		log.Println("request url", r.URL.Path)
		err := json.NewDecoder(r.Body).Decode(&event)
		log.Println("Event:", event)
		if err != nil {
			log.Println("Error decoding event:", err)
			response := Response{"Error decoding event:" + err.Error(), events.Event{}, http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Get the corresponding handler function for the event type
		handlerFunc, ok := Events[event.Type]
		if !ok {
			log.Println("Event type not found:", event.Type)
			response := Response{"Event type not found", events.Event{}, http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		response, err := handlerFunc(event.Payload)
		if err != nil {
			log.Println("Error handling event:", err)
			response := Response{"Error handling event:" + err.Error(), events.Event{}, http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		testData, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling response to JSON:", err)
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		}
		log.Println("Response before sending it :", string(testData))
		json.NewEncoder(w).Encode(response)
	}
}
