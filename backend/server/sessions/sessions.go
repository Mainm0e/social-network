package sessions

import (
	"backend/db"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

/*
getUserID is a helper function for the sessions package, which takes a user's email
address as a string and returns the user's ID as an integer. It is used when a user
logs in to generate a Session struct (which has a userID field). The function also
returns an error value, which is non-nil if an error occurs during the database query.
*/
func getUserID(email string) (int, error) {
	var userID int

	user, err := db.FetchData("users", "email = ?", email)
	userID = user[0].(db.User).UserId
	if err != nil {
		return 0, fmt.Errorf("sessions.getUserID() error: %v", err)
	}

	return userID, nil
}

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

	userID, err := getUserID(username)
	if err != nil {
		return "", errors.New("sessions.Create() error in retrieving userID: " + err.Error())
	}

	session := &Session{
		SessionID: sessionID, // Possibly redundant as it is the key in the sync.Map
		UserID:    userID,
		Admin:     admin,
		Expires:   expires,
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
CookieCheck() checks if the user is authenticated by taking a session cookie
as input and checking if the session ID is present in the SessionStore sync.Map data
structure by calling the *SessionStore Get() method. It returns a boolean indicating
whether the user is authenticated or not, and an error, which is non-nil if an error
occurs during the authentication check.
*/
func CookieCheck(cookie *http.Cookie) (bool, error) {
	_, isValidSession, err := Store.Get(cookie.Value)
	return isValidSession, err
}

/*
SessionCheck() checks if the user is authenticated by taking a session ID as input
and checking if the session ID is present in the SessionStore sync.Map data structure
by calling the *SessionStore Get() method. It returns a boolean indicating whether the
user is authenticated or not, and an error, which is non-nil if an error occurs during
the authentication check.
*/
func SessionCheck(sessionID string) (bool, error) {
	_, isValidSession, err := Store.Get(sessionID)
	return isValidSession, err
}

/*
Login() is a global function in the sessions package, which is used to create a new user
session when a frontend login event has been authenticated. It takes the username and
admin boolean as arguments and calls the *SessionStore Create() method to create a new
user session. It returns the session ID as well as an error, which is non-nil if an error
occurs during the session creation.
*/
func Login(userName string, admin bool) (string, error) {
	sessionID, err := Store.Create(userName, admin)
	if err != nil {
		return "", errors.New("error in sessions.Login(): " + err.Error())
	}
	// Log session creation
	log.Printf("Session created for user \" %s \", with sessionID: %v", userName, sessionID)
	return sessionID, nil
}

/*
Logout() is a global function in the sessions package, which is used to delete a user
session when a frontend logout event has been received. It takes the session ID as an
argument and calls the *SessionStore Delete() method to delete the associated user
session. It returns an error, which is non-nil if an error occurs during the session
deletion.
*/
func Logout(sessionID string) error {
	err := Store.Delete(sessionID)
	if err != nil {
		return errors.New("error in sessions.Logout(): " + err.Error())
	}
	return nil
}
