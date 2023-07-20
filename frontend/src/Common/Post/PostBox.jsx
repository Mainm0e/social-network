import React, { useState, useEffect } from "react";
import Comment from "./CommentBox";
import "./Post.css";
import { getCookie, getUserId } from "../../tools/cookie";
import { checkPostData } from "../../tools/checkdata";
import { fetchData } from "../../tools/fetchData";
import CreateEvent from "../Event/CreateEvent";

const PostList = ({ profileId, groupId, from }) => {
  const [postData, setPostData] = useState(null);
  useEffect(() => {
    const userId = getUserId("userId");
    const sessionId = getCookie("sessionId");
    const method = "POST";
    const type = "getPosts";
    const payload = {
      sessionId: sessionId,
      userId: parseInt(userId),
      from: from,
      profileId: parseInt(profileId),
      groupId: parseInt(groupId),
    };
    fetchData(method, type, payload).then((responseData) => {
      setPostData(responseData);
    });
  }, [profileId, groupId, from]);

  const createPost = () => {
    if (postData !== null) {
      return postData.map((post) => (
        <Post
          key={post.postId}
          id={post.postId}
          title={post.title}
          content={post.content}
          image={post.image}
          time={post.date}
          user={post.creatorProfile}
          comments={post.comments}
        />
      ));
    }
  };

  if (!postData) {
    return null;
  } else {
    return <div className="post_list">{createPost()}</div>;
  }
};
const Post = ({ id, title, content, image, time, user, comments }) => {
  const checkImage = () => {
    if (image === "" || image === null || image === undefined) {
      return null;
    } else {
      return (
        <div className="post_image">
          {" "}
          <img src={image} alt="content" />{" "}
        </div>
      );
    }
  };
  const activePost = (id) => {
    const postList = document.getElementsByClassName("post");
    const activePost = document.querySelector(`[postid="${id}"]`);
    if (id === null) {
      for (let i = 0; i < postList.length; i++) {
        postList[i].classList.remove("hidden");
      }
    } else {
      for (let i = 0; i < postList.length; i++) {
        postList[i].classList.add("hidden");
      }
      activePost.classList.remove("hidden");
    }
  };
  return (
    <>
      <div className="post" postid={id}>
        {checkImage()}
        <div className="post_header">
          <div className="post_header_left">
            <div className="post_header_info">
              <h2>{title}</h2>
              <p>{time}</p>
            </div>
          </div>
          <div className="post_header_right">
            <div className="post_header_user">
              <img src={user.avatar} alt="avatar" />
              <p>
                {user.firstName} {user.lastName}
              </p>
            </div>
          </div>
        </div>
        <div className="post_body">
          <p className="content">{content}</p>
        </div>
        <Comment id={id} comments={comments} activePost={activePost} />
      </div>
    </>
  );
};

const CreatePost = ({ onSubmit, type }) => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [image, setImage] = useState(null);
  const [showImage, setShowImage] = useState(null);
  const [followers, setFollowers] = useState([]);

  const [follower, setFollower] = useState(null);
  useEffect(() => {
    const method = "POST";
    const type = "profileList";
    const payload = {
      sessionId: getCookie("sessionId"),
      userId: getUserId("userId"),
      request: "followers",
    };
    fetchData(method, type, payload).then((data) => {
      setFollower(data);
    });
  }, []);

  const handleTitleChange = (e) => {
    setTitle(e.target.value);
  };

  const handleContentChange = (e) => {
    setContent(e.target.value);
  };

  const handlePrivacyChange = (e) => {
    setPrivacy(e.target.value);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = () => {
      const base64Image = reader.result;
      setImage(base64Image);
      setShowImage(file);
    };
    reader.readAsDataURL(file);
  };

<<<<<<< HEAD
  const handlePrivacyChange = (e) => { 
    setPrivacy(e.target.value);
  };
  
=======
  const handleFollowerChange = (followerId) => {
    if (followers.includes(followerId)) {
      // If follower is already selected, remove it
      setFollowers(followers.filter((follower) => follower !== followerId));
    } else {
      // If follower is not selected, add it
      setFollowers([...followers, followerId]);
    }
  };

>>>>>>> soma
  const handleSubmit = (e) => {
    e.preventDefault();
    const postData = {
      title: title,
      content: content,
      image: image,
      privacy: privacy,
      followers: followers,
    };
    onSubmit(postData);
    setTitle("");
    setContent("");
    setImage(null);
    setFollowers([]);
  };

  return (
    <div className="create_post">
      <form onSubmit={handleSubmit}>
        <div className="create_post_top">
          <input
            type="text"
            placeholder="Title"
            className="create_post_title"
            value={title}
            onChange={handleTitleChange}
            maxLength="20"
          />
        </div>
        <div className="create_post_bottom">
          <textarea
            placeholder="Content"
            className="create_post_content"
            value={content}
            onChange={handleContentChange}
          />
          <input type="file" accept="image/*" onChange={handleImageChange} />
        </div>
        {image && (
          <div className="create_post_image">
            <img src={URL.createObjectURL(image)} alt="Selected" />
          </div>
        )}
        <div className="create_post_button">
          <button type="submit" className="create_post_submit">
            Submit
          </button>
        </div>
<<<<<<< HEAD
=======
        {type === "user" && (
>>>>>>> soma
        <div className="create_post_privacy">
          <select value={privacy} onChange={handlePrivacyChange}>
            <option value="public">Public</option>
            <option value="private">Private</option>
<<<<<<< HEAD
          </select>
        </div>
=======
            <option value="semi-private">Semi-Private</option>
          </select>
            {/* Render the privacy option here */}
            {privacy === "semi-private" && (
              <FollowerList
              users={follower}
              followers={followers}
              handleFollowerChange={handleFollowerChange}
              />
              )}
              </div>
        )}
>>>>>>> soma
      </form>
    </div>
  );
};

const FollowerList = ({ users, followers, handleFollowerChange }) => {
  return (
    <div className="create_post_follower_list">
      {users.map((user) => (
        <div key={user.userId}>
          <label>
            <input
              type="checkbox"
              value={user.userId}
              checked={followers.includes(user.userId)}
              onChange={() => handleFollowerChange(user.userId)}
            />
            {user.firstName}
          </label>
        </div>
      ))}
    </div>
  );
};

// !! Main Component !!
const PostBox = ({ id, from }) => {
  const [body, setBody] = useState("");
  const [pageType, setPageType] = useState("");
  const url = new URL(window.location.href);
  const urlParams = new URLSearchParams(window.location.search);
  const sendId = urlParams.get("id");
  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.substring(1);
      setBody(hash);
    };

    if (url.pathname === "/user") {
      setPageType("user");
    } else if (url.pathname === "/group") {
      setPageType("group");
    }

    // Listen for hash changes in the URL
    window.addEventListener("hashchange", handleHashChange);
    handleHashChange(); // Initialize the body state based on the current hash

    // Cleanup the event listener on unmount
    return () => {
      window.removeEventListener("hashchange", handleHashChange);
    };
  }, []);

  // submit post
  const handleSubmitPost = (postData) => {
    // Logic to handle the submission of the post data
    //check url value
    let groupId = 0;
    if (url.pathname === "/user") {
      groupId = 0;
    } else if (url.pathname === "/group") {
      groupId = parseInt(sendId);
    }

    const check = checkPostData(postData);
    if (check.status === true) {
      const sessionId = getCookie("sessionId");
      // Make API requests or perform other operations here
      // request to create post
      const method = "POST";
      const type = "createPost";
      const payload = {
        sessionId: sessionId,
        postId: 0,
        userId: getUserId("userId"),
        title: postData.title,
        content: postData.content,
        image: postData.image,
        status: postData.privacy,
        groupId: groupId,
        comments: [],
        date: "",
        followers: postData.followers,
      };
      fetchData(method, type, payload).then((data) => {
        window.location.hash = "postlist";
      });
    } else {
      alert(check.message);
    }
  };

  return (
    <>
      {body === "createpost" && (
        <section id="createpost">
          <CreatePost onSubmit={handleSubmitPost} type={pageType} />
        </section>
      )}

      {body === "postlist" && (
        <section id="postlist">
          <PostList
            profileId={from === "profile" ? id : 0}
            groupId={from === "group" ? id : 0}
            from={from}
          />
        </section>
      )}

      {body === "createevent" && (
        <section id="createevent">
          <CreateEvent
            profileId={from === "profile" ? id : 0}
            groupId={from === "group" ? id : 0}
          />
        </section>
      )}
    </>
  );
};

export default PostBox;
