/*
The sessions package is used for general session management and authentication functionality.
It has two primary struct types:

  - Session: stores the session ID and username of the logged in user, an admin boolean to
    indicate if the user is an admin or not, and the expiry date / time of the session.
  - SessionStore: stores a thread-safe map of sessions. It stores the sessions map in a
    sync.Map data structure, which allows for concurrent read / write access to the map.

In addition the sessions package has a number of Global (outward-facing) functions, as well
as methods for the *SessionStore struct mentioned above. The Global functions are the
following:

  - Login: creates a new user session and stores it in the Store.Data sync.Map data
    structure.
  - Logout: deletes the user session from the Store.Data sync.Map data structure.
  - Check: checks if a corresponding valid user session exists in the Store.Data
    sync.Map data structure.
*/
package sessions

import (
	"sync"
	"time"
)

const (
	// Cookie name for user session
	COOKIE_NAME = "supercalafragalisticexpialadoshus"
	// Session duration in seconds
	SESSION_DURATION = 3600 // 1 hour
)

// Store of user sessions
var Store = SessionStore{}

/*
Session struct is used to store the session ID and username of the logged in user,
as well as the expiry date / time of each session. It is used to verify if the user is
logged in, and to clear the session ID when the user logs out or the session expires.
A different variable is assigned to this struct for each user and is added as a key-value
pair to the map in the SessionStore struct. It also has a boolean "Admin" field to
indicate if the user is an admin or not (this is for future use, if admin functionality
is implemented).
*/
type Session struct {
	SessionID string
	UserID    int
	Admin     bool
	Expires   time.Time
}

/*
SessionStore struct is used to store the map of sessions. It is used to initialise the
sessions map from an external calling file (from server package for instance, or the main
file). It is also used to store the sessions map in the sync.Map data structure, which
allows for concurrent read / write access to the map. This is necessary as the sessions
map is accessed by multiple goroutines (multiple users) at the same time primarily via the
server authentication middleware.
*/
type SessionStore struct {
	Data sync.Map
}
