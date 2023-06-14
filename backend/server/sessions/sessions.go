package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

/*
Initialise is a method for the *SessionStore struct, which allows for initialisation of the
sessions map from an external calling file (from server package for instance, or the main file).
*/
func (store *SessionStore) Initialise() {
	store.Data = sync.Map{
		// Initialise the sessions map
	}
	// Must use the following in calling file / function:
	// sessions.Store.Initialise() // Store is a variable of type SessionStore, declared in sessions.go
}

/*
GenerateSessionID is a method for the *SessionStore struct, which generates a unique session
ID for each user session. It is used when a user logs in to generate a session ID for the user
session and store it in the sessions map, as well as return it to the frontend to be set as a
cookie. It checks if the generated session ID already exists in the sessions map (no RWMutex is
used as the *SessionStore struct uses a sync.Map data structure, which is thread-safe).
*/
func (store *SessionStore) GenerateSessionID() (string, error) {
	// Check if *SessionStore struct has been initialised or is present (should have a Data field)
	if store == nil {
		return "", errors.New("error in sessions.<store>.GenerateSessionID(): " +
			"sessions store (sync.Map) not initialised")
	}

	// Loop until a unique session ID is generated
	for {
		sessionID, _ := uuid.NewV4()
		idString := sessionID.String()
		if _, exists := store.Data.Load(idString); !exists {
			return idString, nil
		}
	}
}

/*
Create is a method for the *SessionStore struct, which creates a new user session and stores
it in the SessionStore sync.Map data structure. It is used when a frontend login event is authenticated,
and calls the local helper function GenerateSessionID() to generate a unique session ID for the
user session. It then creates a new Session struct with the generated session ID, username, admin
boolena and expiry date / time, and stores it in the SessionStore. It returns the session ID as well
as an error, which is non-nil if an error occurs during the session creation.
*/
func (store *SessionStore) Create(username string, admin bool) (string, error) {
	sessionID, err := store.GenerateSessionID()
	if err != nil {
		return "", errors.New("error in sessions.<store>.Create(): " + err.Error())
	}

	expires := time.Now().Add(time.Duration(SESSION_DURATION) * time.Second)

	session := &Session{
		ID:       sessionID,
		Username: username,
		Admin:    admin,
		Expires:  expires,
	}

	store.Data.Store(sessionID, session)

	return sessionID, nil
}

/*
CheckAuthentication checks if the user is authenticated by checking the session cookie
in the request. If the user is authenticated, the function returns true. If the user is
not authenticated, the function returns false. The function also returns an error which
is non-nil if an error occurs during the authentication check.
*/
func CheckAuthentication(r *http.Request) (bool, error) {
	cookie, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// if the cookie is not set, return false
			return false, nil
		}
		// For any other type of error, return an error
		return false, err
	}

	// Get the user session from the store
	user, ok := Store.Data.Load(cookie.Value)
	if !ok {
		return false, nil
	}

	// Check if the session is expired
	if cookie.Expires.Before(time.Now()) {
		return false, nil
	}

	return user != "", nil
}
