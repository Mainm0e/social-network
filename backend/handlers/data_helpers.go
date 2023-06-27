package handlers

import (
	"backend/db"
	"backend/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
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
smallProfiles use for followers, followings list in profile page and explore page
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
	var imageUrl string
	if user.Avatar != nil && *user.Avatar != "" {
		imageUrl = *user.Avatar
	} else {
		imageUrl = "./images/avatars/default.png"
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
	followers, err := db.FetchData("follow", "followeeId", userId)
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
func findFollowings(userId int) ([]int, error) {
	followings, err := db.FetchData("follow", "followerId", userId)
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
	if err != nil {
		return Profile{}, errors.New("Error findFollowings: " + err.Error())
	}

	var imageUrl string
	if user.Avatar != nil && *user.Avatar != "" {
		imageUrl = *user.Avatar
	} else {
		imageUrl = "./images/avatars/default.png"
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
login is a function that attempts to log in a user based on the provided data.
It takes in a byte slice `data` containing the login information.
It returns a boolean value indicating whether the login was successful, and an error if any occurred.
*/
func (lg *LoginData) login() (int, error) {

	// Fetch user data from the database based on the provided email.
	user, err := db.FetchData("users", "email", lg.Email)
	if err != nil {
		return 0, errors.New("Error fetching data" + err.Error())
	}

	// Check if a user with the specified email was found.
	if len(user) == 0 {
		return 0, errors.New("user not found")
	}

	// Compare the provided password with the password stored in the database.
	if user[0].(db.User).Password == lg.Password {
		return user[0].(db.User).UserId, nil
	} else {
		return 0, errors.New("password incorrect")
	}
}

/*
register is a function that attempts to register a new user based on the provided data.
It takes in a byte slice `data` containing the registration information.
It returns a boolean value indicating whether the registration was successful, and an error if any occurred.
*/
func (regData *RegisterData) register() error {
	_, err := db.InsertData("users", regData.Email, regData.FirstName, regData.LastName, regData.BirthDate, regData.NickName, regData.Password, regData.AboutMe, regData.Avatar, "public", time.Now())
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
		return errors.New("Error fetching user " + err.Error())
	}
	if len(users) == 0 {
		return errors.New("user not found")
	}
	user := users[0].(db.User)
	// if frontend guys were too lazy to check if privacy changed really or same thing is sent again check it here before updating

	err = db.UpdateData("users", privacy, user.UserId)
	if err != nil {
		return errors.New("Error updating user " + err.Error())
	}
	return nil
}

/*
InsertPost function insert the post into database and check if it is semi-private then it insert the followers that user selected to semiPrivate table
if error occur then it return error
*/
func InsertPost(post Post) error {

	id, err := db.InsertData("posts", post.UserId, post.GroupId, post.Title, post.Content, time.Now(), post.Status, "")
	if err != nil {
		return errors.New("Error inserting post " + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting post ")
	}
	fmt.Println("post image", post.Image)
	if post.Image != "" {
		// Process the image and save it to the local storage
		str := strconv.Itoa(int(id))
		url := "./images/posts/" + str
		url, err := utils.ProcessImage(post.Image, url)
		if err != nil {
			log.Println("Error processing post image:", err)
			//response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return err
		}
		post.Image = url
	} else {
		post.Image = ""
	}
	err = db.UpdateData("posts", post.Image, id)
	if err != nil {
		return errors.New("Error updating post " + err.Error())
	}
	if post.Status == "semi-private" {
		for _, followerId := range post.Followers {
			_, err := db.InsertData("semiPrivate", id, followerId)
			if err != nil {
				return errors.New("Error inserting semiPrivate" + err.Error())
			}
		}
	}
	return nil
}
func InsertComment(comment Comment) error {
	id, err := db.InsertData("comments", comment.UserId, comment.PostId, comment.Content, "", time.Now())
	if err != nil {
		return errors.New("Error inserting comment " + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting comment ")
	}
	if comment.Image != "" {
		// Process the image and save it to the local storage
		str := strconv.Itoa(int(id))
		url := "./images/comments/" + str
		url, err := utils.ProcessImage(comment.Image, url)
		if err != nil {
			log.Println("Error processing comment image:", err)
			//response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return err
		}
		comment.Image = url
	} else {
		comment.Image = ""
	}
	err = db.UpdateData("comments", comment.Image, id)
	if err != nil {
		return errors.New("Error updating comment image" + err.Error())
	}
	return nil
}
func readComments(currentUserId, postId int) ([]Comment, error) {
	/* 	// check if user has permission to see this post
	   	posts, err := db.FetchData("posts", "postId", postId)
	   	if err != nil {
	   		return []Comment{}, errors.New("Error fetching post" + err.Error())
	   	}
	   	if len(posts) == 0 {
	   		return []Comment{}, errors.New("post not found")
	   	}
	   	post := posts[0].(db.Post)
	   	ok, err := checkPost(post, currentUserId)
	   	if err != nil {
	   		return []Comment{}, errors.New("Error checking post" + err.Error())
	   	}
	   	if !ok {
	   		return []Comment{}, errors.New("you don't have permission to see this post")
	   	} */
	// now that we know user has permission to see this post we can fetch comments
	comments, err := db.FetchData("comments", "postId", postId)
	if err != nil {
		return []Comment{}, errors.New("Error fetching comments" + err.Error())
	}
	if len(comments) == 0 {
		return []Comment{}, nil
	}
	var commentsList []Comment
	for _, comment := range comments {
		dbComment := comment.(db.Comment)
		if dbComment.Image != "" {
			image, err := utils.RetrieveImage(dbComment.Image)
			if err != nil {
				return []Comment{}, errors.New("Error retrieving post image: " + err.Error())
			}
			dbComment.Image = image
		}
		CreatorProfile, err := fillSmallProfile(dbComment.UserId)
		if err != nil {
			return []Comment{}, errors.New("Error filling profile" + err.Error())
		}

		commentsList = append(commentsList, Comment{
			CommentId:      dbComment.CommentId,
			PostId:         dbComment.PostId,
			UserId:         dbComment.UserId,
			CreatorProfile: CreatorProfile,
			Content:        dbComment.Content,
			Image:          dbComment.Image,
			Date:           dbComment.CreationTime,
		})

	}
	return commentsList, nil

}

/*
checkPost function check if user has permission to see the post that pass to it base on post status.
it used in readPost and readComments functions.
it return true if user has permission and false if not. if error occur it return error.
*/
func checkPost(dbPost db.Post, userId int) (bool, error) {

	if dbPost.UserId == userId {
		return true, nil
	}
	switch dbPost.Status {
	case "semi-private":
		{
			semiPrivates, err := db.FetchData("semiPrivate", "postId", dbPost.PostId)
			if err != nil {
				return false, errors.New("Error fetching semiPrivate" + err.Error())
			}
			if len(semiPrivates) == 0 {
				return false, errors.New("no followers found")
			}
			for _, semiPrivate := range semiPrivates {
				if semiPrivate.(db.SemiPrivate).UserId == userId {
					return true, nil
				}
			}
		}
	case "private":
		{
			//check if userId is in followers of the post creator
			followers, err := findFollowers(dbPost.UserId)
			if err != nil {
				return false, errors.New("Error fetching followers data" + err.Error())
			}
			for _, follower := range followers {
				if follower == userId {
					return true, nil
				}
			}
		}
	case "group":
		{
			groups, err := db.FetchData("group_memeber", "userId", userId)
			if err != nil {
				return false, errors.New("Error fetching group" + err.Error())
			}
			if len(groups) == 0 {
				return false, errors.New("user is not a memeber of group")
			}
			for _, group := range groups {
				if group.(db.GroupMember).GroupId == dbPost.GroupId {
					return true, nil
				}
			}
		}
	case "public":
		{
			return true, nil
		}
	}
	return false, errors.New("user doesn't have permission to this post")

}

/*
ReadPost function read post from database and check if current user has permission to see it, using checkPost function.
it return post if user has permission, error if error occur or post not found.
*/
func ReadPost(postId int, userId int) (Post, error) {
	dbPosts, err := db.FetchData("posts", "postId", postId)
	if err != nil {
		return Post{}, errors.New("Error fetching post" + err.Error())
	}
	if len(dbPosts) == 0 {
		return Post{}, errors.New("post not found")
	}
	dbPost := dbPosts[0].(db.Post)
	creator, err := fillSmallProfile(dbPost.UserId)
	if err != nil {
		return Post{}, errors.New("Error fetching post creator" + err.Error())
	}
	post := Post{
		PostId:         dbPost.PostId,
		UserId:         dbPost.UserId,
		CreatorProfile: creator,
		Title:          dbPost.Title,
		Content:        dbPost.Content,
		Status:         dbPost.Status,
		GroupId:        dbPost.GroupId,
		Date:           dbPost.CreationTime,
	}
	comments, err := readComments(userId, postId)
	if err != nil {
		return Post{}, errors.New("Error reading comments" + err.Error())
	}
	post.Comments = comments
	post.Followers = []int{}
	if dbPost.Image != "" {
		image, err := utils.RetrieveImage(dbPost.Image)
		if err != nil {
			return Post{}, errors.New("Error retrieving post image: " + err.Error())
		}
		post.Image = image
	}
	ok, err := checkPost(dbPosts[0].(db.Post), userId)

	if err != nil {
		return Post{}, errors.New("Error checking post" + err.Error())
	}
	if !ok {
		return Post{}, errors.New("user doesn't have permission to this post")
	}
	return post, nil

}

/*
ReadPostsByProfile function read all posts of a profile from database
and returns it if user have permission to see it base on post status and user relation to the post creator.
if error occur then it return error.
*/
func ReadPostsByProfile(currentUserId int, userId int) ([]Post, error) {
	var posts []Post
	dbPosts, err := db.FetchData("posts", "userId", userId)
	if err != nil {
		return []Post{}, errors.New("Error fetching posts" + err.Error())
	}
	if len(dbPosts) == 0 {
		return []Post{}, errors.New("posts not found")
	}
	for _, dbPost := range dbPosts {
		post, err := ReadPost(dbPost.(db.Post).PostId, currentUserId)
		if err != nil {
			return []Post{}, errors.New("Error checking post" + err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}

/*
ReadPostsByGroup function read all posts of a group from database
and returns it if user have permission to see it if user is a member of the group
if error occur then it return error.
*/
func ReadPostsByGroup(currentUserId int, groupId int) ([]Post, error) {
	var posts []Post
	dbPosts, err := db.FetchData("posts", "groupId", groupId)
	if err != nil {
		return []Post{}, errors.New("Error fetching posts" + err.Error())
	}
	if len(dbPosts) == 0 {
		return []Post{}, errors.New("posts not found")
	}
	for _, dbPost := range dbPosts {
		post, err := ReadPost(dbPost.(db.Post).PostId, currentUserId)
		if err != nil {
			return []Post{}, errors.New("Error checking post" + err.Error())
		}
		posts = append(posts, post)
	}
	return posts, nil
}

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
func NonMemberUsers(groupId int, userId int, sessionId string) ([]Profile, error) {
	users, err := ReadAllUsers(userId, sessionId)
	if err != nil {
		return []Profile{}, errors.New("Error fetching users: " + err.Error())
	}

	members, err := db.FetchData("group_member", "groupId", groupId)
	if err != nil {
		return []Profile{}, errors.New("Error fetching members: " + err.Error())
	}

	memberIds := make(map[int]struct{})
	for _, member := range members {
		memberIds[member.(db.GroupMember).UserId] = struct{}{}
	}

	var nonMembers []Profile
	for _, user := range users {
		if _, exists := memberIds[user.UserId]; !exists {
			nonMembers = append(nonMembers, user)
		}
	}

	return nonMembers, nil
}

func InsertGroupInvitation(senderId int, groupId int, receiverId int, content string) error {
	_, err := db.InsertData("notifications", receiverId, senderId, groupId, "group_invitation", content, time.Now())
	if err != nil {
		return errors.New("Error inserting group invitation" + err.Error())
	}
	return nil
	// TODO: send notification to receiver
}
func InsertGroupRequest(senderId int, groupId int) error {
	group, err := ReadGroup(groupId)
	if err != nil {
		return errors.New("Error fetching group" + err.Error())
	}
	receiverId := group.CreatorProfile.UserId
	if receiverId == 0 {
		return errors.New("error fetching group creator")
	}
	id, err := db.InsertData("notifications", receiverId, senderId, groupId, "group_request", "", time.Now())
	if err != nil {
		return errors.New("Error inserting group request" + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting group request")
	}
	return nil
}

/*
GroupInvitationCheck function check if user accept or reject or ignore the group invitation, then insert or delete the user from group_member table base on user decision
if error occur then it return error
*/
func GroupInvitationCheck(accept string, notifId int, userId int, groupId int) error {
	if accept == "" {
		return nil
	}
	err := db.DeleteData("notifications", notifId)
	if err != nil {
		return errors.New("Error deleting group invitation" + err.Error())
	}
	if accept == "accept" {
		err := db.UpdateData("group_member", "member", userId)
		if err != nil {
			return errors.New("Error inserting group member" + err.Error())

		} else if accept == "reject" {
			err = db.DeleteData("group_member", userId)
			if err != nil {
				return errors.New("Error deleting group member" + err.Error())
			}
		}
	}
	return nil

}

// todo: change status values in follow table (Maryam)
func InsertFollowRequest(senderId int, receiverId int) error {
	_, err := db.InsertData("notifications", receiverId, senderId, 0, "follow_request", "", time.Now())
	if err != nil {
		return errors.New("Error inserting follow request" + err.Error())
	}
	_, err = db.InsertData("follow", senderId, receiverId, "pending")
	if err != nil {
		return errors.New("Error inserting follow request in follow table" + err.Error())
	}
	return nil
}

func ReadNotifications(userId int, sessionId string) ([]Notification, error) {
	notifications, err := db.FetchData("notifications", "receiverId", userId)
	if err != nil {
		return []Notification{}, errors.New("Error fetching notifications" + err.Error())
	}
	result := make([]Notification, len(notifications))
	for i, n := range notifications {
		if notification, ok := n.(db.Notification); ok {
			profile, err := fillSmallProfile(notification.SenderId)
			if err != nil {
				return []Notification{}, errors.New("Error fetching profile" + err.Error())
			}
			result[i] = Notification{
				SessionId:    sessionId,
				Notification: notification,
				Profile:      profile,
			}
		} else {
			return nil, fmt.Errorf("invalid notification type at index %d", i)
		}
	}

	return result, nil
}

/*
DeleteFollowRequest function delete the follow request from notification table and update the follow table base on user decision
if error occur then it return error
*/

func DeleteFollowRequest(followId int, notifId int, response string) error {
	err := db.DeleteData("notifications", followId)
	if err != nil {
		return errors.New("Error deleting follow request" + err.Error())
	}
	if response == "accept" {
		err = db.UpdateData("follow", "follower", followId)
		if err != nil {
			return errors.New("Error updating follow request" + err.Error())
		}
	} else if response == "reject" {
		err = db.DeleteData("follow", followId)
		if err != nil {
			return errors.New("Error deleting follow request" + err.Error())
		}
	}
	return nil
}
