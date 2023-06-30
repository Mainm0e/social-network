package handlers

import (
	"backend/db"
	"backend/events"
	"backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

/*
Here we have the handlers for the post requests
- CreatePost: creates a new post in the database and returns the sessionId
- GetPost: gets a post from the database and returns the post and the sessionId
- GetPosts: gets all posts related to a user from the database and returns the posts and the sessionId
*/
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
	err = InsertPost(post)
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

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func readComments(currentUserId, postId int) ([]Comment, error) {
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
	fmt.Println("status", dbPost.Status)
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
func GetPost(payload json.RawMessage) (Response, error) {
	var response Response
	var request RequestPost
	err := json.Unmarshal(payload, &request)
	log.Println("User: ", request)
	if err != nil {
		// handle the error
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.SessionId == "" {
		response = Response{"sessionId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	if request.PostId == 0 {
		response = Response{"postId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}
	// TODO: check if the userId is necessary to get from request
	if request.UserId == 0 {
		response = Response{"userId is required", events.Event{}, http.StatusBadRequest}
		return response, err
	}

	//get post from database
	post, err := ReadPost(request.PostId, request.UserId)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	post.SessionId = request.SessionId
	payload, err = json.Marshal(post)
	if err != nil {
		response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
		return response, err
	}
	event := events.Event{
		Type:    "requestPost",
		Payload: payload,
	}
	return Response{"post retrieved successfully!", event, http.StatusOK}, nil
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
	log.Println("User: ", request)
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
		posts, err = ReadPostsByGroup(request.UserId, request.GroupId)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else if request.From == "profile" {
		posts, err = ReadPostsByProfile(request.UserId, request.ProfileId)
		if err != nil {
			response = Response{err.Error(), events.Event{}, http.StatusBadRequest}
			return response, err
		}
	} else {
		response = Response{"from is required", events.Event{}, http.StatusBadRequest}
		return response, err
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
