package handlers

import (
	"backend/db"
	"backend/utils"
	"errors"
)

/*
fetchUser get a user from the database using the userId.
It returns the user and any error encountered during the process.
*/
func fetchUser(key string, value any) (db.User, error) {
	users, err := db.FetchData("users", key+" = ?", value)
	if err != nil {
		return db.User{}, errors.New("Error fetching user data: " + err.Error())
	}
	if len(users) == 0 {
		return db.User{}, errors.New("user not found")
	}
	return users[0].(db.User), nil
}

/*
smallProfiles use for followers, followings list in profile page and explore page
*/

/*
fillSmallProfile fills a SmallProfile struct with data from a db.User struct.
It returns the filled SmallProfile struct and any error encountered during the process.
*/
func fillSmallProfile(userId int) (SmallProfile, error) {
	user, err := fetchUser("userId", userId)
	if err != nil {
		return SmallProfile{}, errors.New("Error fetchingUser:" + err.Error())
	}
	smallProfile := SmallProfile{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
	var imageUrl string
	if user.Avatar != nil && *user.Avatar != "" {
		imageUrl = *user.Avatar
	} else {
		imageUrl = "./images/avatars/default.gif"
	}
	avatar, err := utils.RetrieveImage(imageUrl)
	if err != nil {
		return SmallProfile{}, errors.New("Error retrieving avatar image: " + err.Error())
	}

	smallProfile.Avatar = &avatar
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
	followers, err := db.FetchData("follow", "followeeId = ?", userId)
	if len(followers) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("Error fetching followers data" + err.Error())
	}
	var followerIds []int
	for _, follower := range followers {
		if follower.(db.Follow).Status == "following" {
			followerIds = append(followerIds, follower.(db.Follow).FollowerId)
		}
	}
	return followerIds, nil
}

/*
findFollowings() takes a userID integer and returns a slice of integers containing the IDs of the users that the user with the given ID is following.
It also returns an error value, which is non-nil if an error occurred during the process.
*/
func findFollowings(userId int) ([]int, error) {
	followings, err := db.FetchData("follow", "followerId = ?", userId)
	if len(followings) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("Error fetching followings data" + err.Error())
	}
	var followingIds []int
	for _, following := range followings {
		if following.(db.Follow).Status == "following" {
			followingIds = append(followingIds, following.(db.Follow).FolloweeId)
		}
	}
	return followingIds, nil
}

/*
this function check the current user and requested user profile relation
if current user follow requested user then it return 'following' if they already request to follow then it return 'pending' else it return 'follow'
if error occur then it return error
*/
func checkUserRelation(userId int, profileId int) (string, error) {
	followings, err := db.FetchData("follow", "followerId = ?", userId)
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
	user, err := fetchUser("userId", profileId)
	if err != nil {
		return Profile{}, errors.New("Error fetchingUser:" + err.Error())
	}
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
	if err != nil {
		return Profile{}, errors.New("Error findFollowings: " + err.Error())
	}

	var imageUrl string
	if user.Avatar != nil && *user.Avatar != "" {
		imageUrl = *user.Avatar
	} else {
		imageUrl = "./images/avatars/default.gif"
	}
	avatar, err := utils.RetrieveImage(imageUrl)
	if err != nil {
		return Profile{}, errors.New("Error retrieving avatar image: " + err.Error())
	}
	user.Avatar = &avatar

	profile := Profile{
		sessionId,
		user.UserId,
		user.NickName.String,
		user.FirstName,
		user.LastName,
		user.Avatar,
		"",
		user.Privacy,
		len(followers),
		len(followings),
		PrivateProfile{},
	}
	if user.UserId == userId {
		profile.PrivateData = PrivateProfile{
			user.BirthDate,
			user.Email,
			user.AboutMe.String,
			followers,
			followings,
		}
		profile.Relation = "you"
		return profile, nil

	}
	profile.Relation = status
	if status == "following" || user.Privacy == "public" {
		profile.PrivateData = PrivateProfile{
			user.BirthDate,
			user.Email,
			user.AboutMe.String,
			followers,
			followings,
		}
		return profile, nil
	} else if status == "pending" || status == "follow" {
		return profile, nil
	} else {
		return Profile{}, errors.New("error checkUserRelation: wtf:" + status)
	}
}

/*
ReadAllUsers function return all users in database except current user
for each user it call FillProfile function to fill the profile struct
base on the relation between current user and requested user profile
if error occur then it return error else it returns profile struct and nil.
*/
func ReadAllUsers(userId int, sessionId string) ([]Profile, error) {
	dbUsers, err := db.FetchData("users", "")
	if err != nil {
		return []Profile{}, errors.New("Error fetching users" + err.Error())
	}
	if len(dbUsers) == 0 {
		return []Profile{}, errors.New("no user found")
	}
	var users []Profile
	for _, dbUser := range dbUsers {
		dbUser := dbUser.(db.User)
		user, err := FillProfile(userId, dbUser.UserId, sessionId)
		if err != nil {
			return []Profile{}, errors.New("Error fetching user" + err.Error())
		}
		if user.UserId != userId {
			users = append(users, user)
		}
	}
	return users, nil
}

/*
GetAllGroupMemberIDs returns all the user ids of the members of a group.
It takes a groupID integer as an argument and returns a slice of integers
and an error value, which is non-nil if any of the database operations
failed.
*/
func GetAllGroupMemberIDs(groupId int) ([]int, error) {
	var userIds []int

	// Fetch all group members from the database
	dbGroupMembers, err := db.FetchData("group_member", "groupId = ?", groupId)
	if err != nil {
		return userIds, errors.New("Error fetching group members" + err.Error())
	}

	// If no group members were found, return an error
	if len(dbGroupMembers) == 0 {
		return userIds, errors.New("no group member found")
	}

	// Iterate over all group members and append their user ids to the slice
	for _, dbGroupMember := range dbGroupMembers {
		dbGroupMember := dbGroupMember.(db.GroupMember)
		userIds = append(userIds, dbGroupMember.UserId)
	}

	return userIds, nil
}
