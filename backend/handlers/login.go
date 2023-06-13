package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var loginData LoginData
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			log.Println("Error decoding login data:", err)
			response := Response{false, err.Error(), http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			return
		}
		log.Println("Login try by:", loginData)

		_, err = loginData.login()
		if err != nil {
			response := Response{false, err.Error(), http.StatusBadRequest}
			json.NewEncoder(w).Encode(response)
			log.Println("Error logging in:", err)
			return
		} else {
			response := Response{true, "Login approved", 200}
			json.NewEncoder(w).Encode(response)
			log.Println("Login approved", loginData)
		}
	}
}
