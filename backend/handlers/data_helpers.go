package handlers

import (
	"backend/db"
	"errors"
	"log"
	"time"
)

/*
findFollowers and findFollowings function return the list of followers and followings of the desired user
if error occur then it return error
*/
func findFollowers(userId int) ([]int, error) {
	followers, err := db.FetchData("follow", "followeeId", userId)
	if err != nil {
		return nil, errors.New("Error fetching followers data" + err.Error())
	}
	var followerIds []int
	for _, follower := range followers {
		followerIds = append(followerIds, follower.(db.Follow).FollowerId)
	}
	return followerIds, nil
}
func findFollowings(userId int) ([]int, error) {
	followings, err := db.FetchData("follow", "followerId", userId)
	if err != nil {
		return nil, errors.New("Error fetching followings data" + err.Error())
	}
	var followingIds []int
	for _, following := range followings {
		followingIds = append(followingIds, following.(db.Follow).FolloweeId)
	}
	return followingIds, nil
}

/*
this function check the online user and requested user profile relation
if online user follow requested user then it return 'following' if they already request to follow then it return 'pending' else it return 'follow'
if error occur then it return error
*/
func checkUserRelation(userId int, profileId int) (string, error) {
	followings, err := db.FetchData("follow", "followerId", userId)
	if err != nil {
		return "", errors.New("Error fetching follow data" + err.Error())
	}
	if len(followings) == 0 {
		return "follow", nil
	}
	for _, following := range followings {
		if following.(db.Follow).FolloweeId == profileId {
			return following.(db.Follow).Status, nil
		}
	}
	return "follow", nil
}
func FillProfile(userId int, profileId int) (Profile, error) {

	users, err := db.FetchData("users", "userId", profileId)
	if err != nil {
		return Profile{}, errors.New("Error fetching user data" + err.Error())
	}
	if len(users) == 0 {
		return Profile{}, errors.New("user not found")
	}
	user := users[0].(db.User)
	// check relation between online user and requested user
	status, err := checkUserRelation(userId, profileId)
	if err != nil {
		return Profile{}, errors.New("Error CheckUserRelation: " + err.Error())
	}
	//found number of followers and followings of requested user
	followers, err := findFollowers(profileId)
	if err != nil {
		return Profile{}, errors.New("Error findFollowers: " + err.Error())
	}
	followings, err := findFollowings(profileId)
	log.Println("followings", followings, "followers", followers)
	if err != nil {
		return Profile{}, errors.New("Error findFollowings: " + err.Error())
	}
	profile := Profile{
		user.UserId,
		user.NickName.String,
		user.FirstName,
		user.LastName,
		user.Avatar,
		len(followers),
		len(followings),
		PrivateProfile{},
	}

	if status == "following" || userId == profileId {
		log.Println("yay user is user", profile.PrivateData)
		profile.PrivateData = PrivateProfile{
			user.BirthDate,
			user.Email,
			user.AboutMe.String,
			followers,
			followings,
		}
		log.Println("yay user is user", profile.PrivateData)
	} else if status == "follow" || status == "pending" {
		return profile, nil
	} else {
		return Profile{}, errors.New("error checkUserRelation: wtf")
	}
	return profile, nil
}

/*
login is a function that attempts to log in a user based on the provided data.
It takes in a byte slice `data` containing the login information.
It returns a boolean value indicating whether the login was successful, and an error if any occurred.
*/
func (lg *LoginData) login() (bool, error) {

	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", lg.Email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email was found.
	if len(user) == 0 {
		return false, errors.New("user not found")
	}

	// Compare the provided password with the password stored in the database.
	if user[0].(db.User).Password == lg.Password {
		return true, nil
	} else {
		return false, errors.New("password incorrect")
	}
}

/*
register is a function that attempts to register a new user based on the provided data.
It takes in a byte slice `data` containing the registration information.
It returns a boolean value indicating whether the registration was successful, and an error if any occurred.
*/
func (regData *RegisterData) register() error {
	_, err := db.InsertData("users", regData.NickName, regData.FirstName, regData.LastName, regData.BirthDate, regData.Email, regData.Password, regData.AboutMe, regData.Avatar, "public", time.Now())
	if err != nil {
		return errors.New("Error inserting user" + err.Error())
	}
	return nil
}

/*
IsNotUser is a function that checks if a user with the specified email already exists.
It takes in a string `email` containing the email of the user to check.
It returns a boolean value indicating whether the user exists, and an error if any occurred.
*/
func IsNotUser(email string) (bool, error) {
	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", email)
	if err != nil {
		return false, errors.New("Error fetching data" + err.Error())
	}
	// Check if a user with the specified email already exists.
	if len(user) == 0 {
		// Insert the new user data into the database.
		return true, nil
	} else {
		return false, nil
	}
}

/*
func updateProfile(email string,data any) error {
} */