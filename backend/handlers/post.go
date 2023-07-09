package handlers

import (
	"backend/db"
	"backend/events"
	"backend/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
Here we have the handlers for the post requests
- CreatePost: creates a new post in the database and returns the sessionId
- GetPosts: gets all posts related to a user from the database and returns the posts and the sessionId
*/

/*
readComments function read all comments related to a post and return it as a list of comments. if error occur it return error.
*/
func readComments(postId int) ([]Comment, error) {
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
func checkPostPermission(dbPost db.Post, userId int) (bool, error) {

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
	// if user doesn't have permission to see the post
	return false, nil
}

/*
ReadPost function read post from database and check if current user has permission to see it, using checkPost function.
it return post if user has permission, error if error occur or post not found.
*/
func readPost(postId int, userId int) (Post, error) {
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
	comments, err := readComments(postId)
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
	ok, err := checkPostPermission(dbPosts[0].(db.Post), userId)

	if err != nil {
		return Post{}, errors.New("Error checking post" + err.Error())
	}
	if !ok {
		// if user doesn't have permission to see the post
		return Post{}, nil
	}
	return post, nil

}

/*
readPosts function read all posts base on key that's the column name and value that's the value of the column(posts could be group posts or user posts).
it use readPost function to read each post and check if user has permission to see it.
it return list of posts if found, error if error occur or posts not found.
*/
func readPosts(currentUserId int, key string, value int, groupFlag bool) ([]Post, error) {
	var posts []Post
	dbPosts, err := db.FetchData("posts", key, value)
	if err != nil {
		return []Post{}, errors.New("Error fetching posts" + err.Error())
	}
	if len(dbPosts) == 0 {
		return []Post{}, nil
	}
	for _, dbPost := range dbPosts {
		if groupFlag && dbPost.(db.Post).GroupId == 0 || !groupFlag && dbPost.(db.Post).GroupId != 0 {
			continue
		}

		post, err := readPost(dbPost.(db.Post).PostId, currentUserId)
		if err != nil {
			return []Post{}, errors.New("Error checking post" + err.Error())
		}
		if post.PostId != 0 {
			posts = append(posts, post)
		}
	}
	return posts, nil
}
func GetPosts(payload json.RawMessage) (Response, error) {
	var response Response
	var request ReqAllPosts
	err := json.Unmarshal(payload, &request)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//get posts from database
	var posts []Post
	if request.From == "group" {
		posts, err = readPosts(request.UserId, "groupId", request.GroupId, true)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else if request.From == "profile" {
		posts, err = readPosts(request.UserId, "userId", request.ProfileId, false)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else if request.From == "home" {
		posts, err = readPosts(request.UserId, "userId", request.UserId, false)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else {
		response = Response{"from is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if len(posts) == 0 {

		response = Response{"no posts found", events.Event{}, http.StatusOK}
		return response, nil
	}
	payload, err = json.Marshal(posts)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "requestPosts",
		Payload: payload,
	}
	return Response{"posts retrieved successfully!", event, http.StatusOK}, nil
}

/*
InsertPost function insert the post into database and check if it is semi-private then it insert the followers that user selected to semiPrivate table
if error occur then it return error
*/
func insertPost(post Post) error {

	id, err := db.InsertData("posts", post.UserId, post.GroupId, post.Title, post.Content, time.Now(), post.Status, "")
	if err != nil {
		return errors.New("Error inserting post " + err.Error())
	}
	if id == 0 {
		return errors.New("error inserting post ")
	}
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

func CreatePost(payload json.RawMessage) (Response, error) {
	var response Response
	var post Post
	err := json.Unmarshal(payload, &post)
	log.Println("User: ", post)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}

	//insert new post into database
	err = insertPost(post)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	//send back sessionId
	payload, err = json.Marshal(map[string]string{"sessionId": post.SessionId})
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "createPost",
		Payload: payload,
	}

	response = Response{"post created successfully!", event, http.StatusOK}
	return response, nil
}
