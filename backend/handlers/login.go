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
			return
		}
		log.Println("Login try by:", loginData)

		approve, err := loginData.login()
		if err != nil {
			log.Println("Error logging in:", err)
		}
		if approve {
			response := Response{true, "Login approved"}
			json.NewEncoder(w).Encode(response)
			log.Println("Login approved", loginData)
		} else {
			response := Response{false, "Login denied"}
			json.NewEncoder(w).Encode(response)
		}
	}
}
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("RegisterPage...")
		var registerData RegisterData
		err := json.NewDecoder(r.Body).Decode(&registerData)
		if err != nil {
			log.Println("Error decoding register data:", err)
			return
		}
		log.Println("trying to Register:", registerData)
		approve, err := registerData.register()
		if err != nil {
			log.Println("Error registering:", err)
		}
		if approve {
			response := Response{true, "Registration approved"}
			json.NewEncoder(w).Encode(response)
			log.Println("Registration approved", registerData)
		} else {
			response := Response{false, "Registration denied"}
			json.NewEncoder(w).Encode(response)
		}
	}
}
