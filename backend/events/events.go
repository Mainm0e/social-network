package events

import (
	"encoding/json"
	"fmt"
)

/*
PackageErrorEvent takes an ErrorMessage struct as input and returns a pointer to an
Event struct, as well as an error value, which is non-nil if there was an error encountered
during the marshalling of the ErrorMessage struct. It is used to package an ErrorMessage
struct into an Event struct, which can then be converted to a JSON byte slice and sent
to the frontend / client.
*/
func PackageErrorEvent(errMsg ErrorMessage) (*Event, error) {
	errMsgBytes, err := json.Marshal(errMsg)
	if err != nil {
		return nil, fmt.Errorf("PackageErrorEvent() error - failed to marshal ErrorMessage: %v", err)
	}
	return &Event{
		Type:    "error",
		Payload: errMsgBytes,
	}, nil
}
