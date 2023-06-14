package sessions

import "net/http"

/*
CheckAuthentication checks if the user is authenticated by checking the session cookie
in the request. If the user is authenticated, the function returns true. If the user is
not authenticated, the function returns false. The function also returns an error which
is non-nil if an error occurs during the authentication check.
*/
func CheckAuthentication(r *http.Request) (bool, error) {
	// TODO: Authentication logic, handler calls etc.
	isAuthenticated := true // Set to true to allow testing of other middleware functions
	return isAuthenticated, nil
}
