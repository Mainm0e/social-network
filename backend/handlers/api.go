package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		log.Println("Event:", event)
		if err != nil {
			log.Println("Error decoding event:", err)
			response := Response{false, err.Error(), http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Get the corresponding handler function for the event type
		handlerFunc, ok := Events[event.Event_type]
		if !ok {
			log.Println("Event type not found:", event.Event_type)
			response := Response{false, "Event type not found", http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		response, err := handlerFunc(event.Payload)
		if err != nil {
			log.Println("Error handling event:", err)
			response := Response{false, err.Error(), http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(response)
	}
}
