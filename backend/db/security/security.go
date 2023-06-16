package security

import (
	"golang.org/x/crypto/bcrypt"
)

/*
HashPwd takes a password as a input slice of bytes and uses the bcrypt
package to generate a hashed version of the password which is returned
as a string, along with an error value which is non-nil in the event that
the call to bcrypt's GenerateFromPassword returns an error. Bcrypt uses
hash iterations instead of 'salting', whereby the concept of 'cost' is
used (an integer value for the number of hash iterations applied).
Hashing time is calculated as 2*cost, where variables bcrypt.MinCost is 4
and bcrypt.MaxCost is 31.
*/
func HashPwd(pwd []byte, hashCost int) (string, error) {
	// Manually revert to 'min' or 'max' costs instead of bcrypt.DefaultCost of 10
	if hashCost < bcrypt.MinCost {
		hashCost = bcrypt.MinCost
	} else if hashCost > bcrypt.MaxCost {
		hashCost = bcrypt.MaxCost
	}
	hash, err := bcrypt.GenerateFromPassword(pwd, hashCost)
	if err != nil {
		return "", err
	}
	// Convert resulting byte slice as a string and return it
	return string(hash), nil
}

/*
MatchPasswords takes an input non-hashed password as a slice of bytes, as well
as a hashed password string. It then calls bcrypt's CompareHashAndPassword function
and does a direct conditional compare on the returned error, returning true if
CompareHashAndPassword returns a nil error, otherwise returning false.
*/
func MatchPasswords(plainPwd []byte, hashedPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	// Return false in the event err != nil, else return true
	return err == nil
}
