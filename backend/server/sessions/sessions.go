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
Exists is a method for the *SessionStore struct, which checks if the *SessionStore struct has
been initialised or is present (should have a Data field). It returns an error which is non-nil
if the *SessionStore struct has not been initialised or is not present.
*/
func (store *SessionStore) Exists() error {
	if store == nil {
		return errors.New("error in sessions.<store>.Exists(): session store not initialised (no Data field)")
	}
	return nil
}

/*
GenerateSessionID is a method for the *SessionStore struct, which generates a unique session
ID for each user session. It is used when a user logs in to generate a session ID for the user
session and store it in the sessions map, as well as return it to the frontend to be set as a
cookie. It checks if the generated session ID already exists in the sessions map (no RWMutex is
used as the *SessionStore struct uses a sync.Map data structure, which is thread-safe).
*/
func (store *SessionStore) GenerateSessionID() (string, error) {
	// SessionStore integrity check
	if err := store.Exists(); err != nil {
		return "", errors.New("error in sessions.<store>.GenerateSessionID(): " + err.Error())
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
	// SessionsStore integrity check
	if err := store.Exists(); err != nil {
		return "", errors.New("error in sessions.<store>.Create(): " + err.Error())
	}

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
Delete is a method for the *SessionStore struct, which deletes a user session from the
SessionStore sync.Map data structure. It is used when a frontend logout event is received.
It takes the session ID as an argument and deletes the session from the SessionStore, regardless
of whether the session has expired or not.
*/
func (store *SessionStore) Delete(sessionID string) error {
	// SessionsStore integrity check
	if err := store.Exists(); err != nil {
		return errors.New("error in sessions.<store>.Delete(): " + err.Error())
	}

	store.Data.Delete(sessionID)
	return nil
}

/*
Get is a method for the *SessionStore struct, which gets a user session from the SessionStore
sync.Map data structure. It is used when the session relating to a frontend request needs to be
retrieved. It takes the session ID as an argument and returns the associated Session struct, a
boolean indicating whether the session was found or not, and an error, which is non-nil if an
error occurs during the session retrieval.
*/
func (store *SessionStore) Get(sessionID string) (*Session, bool, error) {
	// SessionsStore integrity check
	if err := store.Exists(); err != nil {
		return nil, false, errors.New("error in sessions.<store>.Get(): " +
			err.Error())
	}

	// Load the session from the store
	value, found := store.Data.Load(sessionID)
	if !found {
		return nil, false, errors.New("error in sessions.<store>.Get(): " +
			"session not found")
	}

	// Assert the value is of type *Session, and check if the session has expired
	session, ok := value.(*Session)
	if !ok || session.Expires.Before(time.Now()) {
		return nil, false, errors.New("error in sessions.<store>.Get(): " +
			"session invalid or has expired")
	}

	return session, true, nil
}

/*
CheckAuthentication checks if the user is authenticated by taking a session cookie
as input and checking if the session ID is present in the SessionStore sync.Map data
structure by calling the *SessionStore Get() method. It returns a boolean indicating
whether the user is authenticated or not, and an error, which is non-nil if an error
occurs during the authentication check.
*/
func CheckAuthentication(cookie *http.Cookie) (bool, error) {
	_, isValidSession, err := Store.Get(cookie.Value)
	return isValidSession, err
}

/*
Login is a global function in the sessions package, which is used to create a new user
session when a frontend login event has been authenticated. It takes the username and
admin boolean as arguments and calls the *SessionStore Create() method to create a new
user session. It returns the session ID as well as an error, which is non-nil if an error
occurs during the session creation.
*/
func Login(userName string, admin bool) (string, error) {
	return Store.Create(userName, admin)
}
