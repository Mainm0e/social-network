package handlers

import (
	"backend/db"
	"backend/utils"
	"errors"
	"fmt"
	"log"
	"time"
)

/*
fetchUser get a user from the database using the userId.
It returns the user and any error encountered during the process.
*/
func fetchUser(userId int) (db.User, error) {
	users, err := db.FetchData("users", "userId", userId)
	if err != nil {
		return db.User{}, errors.New("Error fetching user data" + err.Error())
	}
	if len(users) == 0 {
		return db.User{}, errors.New("user not found")
	}
	return users[0].(db.User), nil
}

/*
smallProfiles use for followers and followings list in profile page and maybe explore page in future
*/

/*
fillSmallProfile fills a SmallProfile struct with data from a db.User struct.
It returns the filled SmallProfile struct and any error encountered during the process.
*/
func fillSmallProfile(userId int) (SmallProfile, error) {
	user, err := fetchUser(userId)
	if err != nil {
		return SmallProfile{}, errors.New("Error fetchingUser:" + err.Error())
	}
	smallProfile := SmallProfile{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
	if user.Avatar != nil {
		avatar, err := utils.RetrieveAvatarImage(*user.Avatar)
		if err != nil {
			return SmallProfile{}, errors.New("Error retrieving avatar image: " + err.Error())
		}
		smallProfile.Avatar = &avatar
	}
	return smallProfile, nil
}

func SmallProfileList(userId int, listName string) ([]SmallProfile, error) {
	var list []int
	var err error
	if listName == "followers" {
		list, err = findFollowers(userId)
		if err != nil {
			return nil, errors.New("Error fetching followers data" + err.Error())
		}
	} else if listName == "followings" {
		list, err = findFollowings(userId)
		if err != nil {
			return nil, errors.New("Error fetching followings data" + err.Error())
		}
	} else {
		//TODO: if we need a small profile list for other things like explore page we should add it here
		return nil, errors.New("invalid list name")
	}
	var smallProfiles []SmallProfile
	for _, followerId := range list {
		smallProfile, err := fillSmallProfile(followerId)
		if err != nil {
			return nil, errors.New("Error fetching followers data" + err.Error())
		}
		smallProfiles = append(smallProfiles, smallProfile)
	}
	return smallProfiles, nil
}

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
this function check the current user and requested user profile relation
if current user follow requested user then it return 'following' if they already request to follow then it return 'pending' else it return 'follow'
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

/*
FillProfile function fill the profile struct with the data base on the relation between current user and requested user profile
if error occur then it return error else it return profile struct and nil.
*/
func FillProfile(userId int, profileId int, sessionId string) (Profile, error) {

	users, err := db.FetchData("users", "userId", profileId)
	if err != nil {
		return Profile{}, errors.New("Error fetching user data" + err.Error())
	}
	if len(users) == 0 {
		return Profile{}, errors.New("user not found")
	}
	user := users[0].(db.User)
	// check relation between current user and requested user
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
	if user.Avatar != nil {

		avatar, err := utils.RetrieveAvatarImage(*user.Avatar)
		if err != nil {
			return Profile{}, errors.New("Error retrieving avatar image: " + err.Error())
		}
		user.Avatar = &avatar
	}
	profile := Profile{
		sessionId,
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
	_, err := db.InsertData("users", regData.Email, regData.FirstName, regData.LastName, regData.BirthDate, regData.NickName, regData.Password, regData.AboutMe, regData.Avatar, "public", time.Now())
	fmt.Println("regData", regData)
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
UpdateProfile is a function that updates the privacy of a user with the specified email.
returns error if any occurred.
*/
func UpdateProfile(email string, privacy string) error {
	users, err := db.FetchData("users", "email", email)
	if err != nil {
		return errors.New("Error fetching user" + err.Error())
	}
	if len(users) == 0 {
		return errors.New("user not found")
	}
	user := users[0].(db.User)
	// if frontend guys were too lazy to check if privacy changed really or same thing is sent again check it here before updating

	err = db.UpdateData("users", privacy, user.UserId)
	if err != nil {
		return errors.New("Error updating user" + err.Error())
	}
	return nil
}

/*
 */
func InsertPost(post Post) error {
	id, err := db.InsertData("posts", post.UserId, post.Title, post.Content, time.Now(), post.Status, post.GroupId)
	if err != nil {
		return errors.New("Error inserting post" + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting post")
	}
	if post.Status == "semi-private" {
		for followerId := range post.Followers {
			_, err := db.InsertData("semiPrivate", id, followerId)
			if err != nil {
				return errors.New("Error inserting semiPrivate" + err.Error())
			}
		}
	}
	return nil
}

/* func GetPost(userId int) (Post, error) {
	dbPosts, err := db.FetchData("posts", "userId", userId)
	if err != nil {
		return Post{}, errors.New("Error fetching post" + err.Error())
	}
	if len(dbPosts) == 0 {
		return Post{}, errors.New("user doesn't have any post")
	}
	var posts []Post
	for _, post := range dbPosts {
		dbPost := post.(db.Post)
		posts = append(posts, Post{
			Post.PostId:  dbPost.PostId,
			Post.UserId:  dbPost.UserId,
			Post.Title:   dbPost.Title,
			Post.Content: dbPost.Content,
			Post.Status:  dbPost.Status,
			Post.GroupId: dbPost.GroupId,
			Post.Date:    dbPost.Date,
		})
		if dbPost.Status == "semi-private" {

		}
	}
}
*/
