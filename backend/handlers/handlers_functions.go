package handlers

import (
	"encoding/json"
	"errors"
)

func login(data []byte) error {
	var login LoginData
	err := json.Unmarshal(data, &login)
	if err != nil {
		return errors.New("Error unmarshalling data" + err.Error())
	}

	return nil
}
