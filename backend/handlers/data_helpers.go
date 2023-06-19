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
		avatar, err := utils.RetrieveImage(*user.Avatar)
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

		avatar, err := utils.RetrieveImage(*user.Avatar)
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
		commentsList = append(commentsList, Comment{
			CommentId: dbComment.CommentId,
			PostId:    dbComment.PostId,
			UserId:    dbComment.UserId,
			Content:   dbComment.Content,
			Date:      dbComment.CreationTime,
		})
	}
	return commentsList, nil

}

/*
ReadPost function read the post from database and return it if user have permission to see it base on
post status and user relation to the post creator.
if error occur then it return error
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
	fmt.Println("post image is: ", dbPost.Image)
	if dbPost.Image != "" {
		image, err := utils.RetrieveImage(dbPost.Image)
		if err != nil {
			return Post{}, errors.New("Error retrieving post image: " + err.Error())
		}
		post.Image = image
		fmt.Println("post image is: ", post.Image)
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

//func ProfilePost(currentUserId)
