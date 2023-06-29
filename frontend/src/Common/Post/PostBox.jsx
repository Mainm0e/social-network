import React, { useState, useEffect } from "react";
import Comment from "./CommentBox";
import "./Post.css";
import { getCookie, getUserId } from "../../tools/cookie";
import { checkPostData } from "../../tools/checkdata";
import { fetchData } from "../../tools/fetchData";

const PostList = ({ profileId }) => {
  const [postData, setPostData] = useState(null);
  useEffect(() => {
    const userId = getUserId("userId");
    const sessionId = getCookie("sessionId");
    const method = "POST";
    const type = "GetPosts";
    const payload = {
      sessionId: sessionId,
      userId: parseInt(userId),
      from: "profile",
      profileId: parseInt(profileId),
      groupId: 0,
    };
    fetchData(method, type, payload).then((responseData) => {
      setPostData(responseData);
    });
  }, []);

  const createPost = () => {
    if (postData !== null) {
      console.log("in createPost", postData.length);
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
    return <div>Loading...</div>;
  } else {
    return <div className="post_list">{createPost()}</div>;
  }
};
const Post = ({ id, title, content, image, time, user, comments }) => {
  const postId = id;

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

const CreatePost = ({ onSubmit }) => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [image, setImage] = useState(null);
  const [showImage, setShowImage] = useState(null);

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
  const handleSubmit = (e) => {
    e.preventDefault();
    const postData = {
      title: title,
      content: content,
      image: image,
      privacy: privacy,
    };
    onSubmit(postData);
    setTitle("");
    setContent("");
    setImage(null);
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
        {showImage && (
          <div className="create_post_image">
            <img src={URL.createObjectURL(showImage)} alt="Selected" />
          </div>
        )}
        <div className="create_post_button">
          <button type="submit" className="create_post_submit">
            Submit
          </button>
        </div>
        <div className="create_post_privacy">
          <select value={privacy} onChange={handlePrivacyChange}>
            <option value="public">Public</option>
            <option value="private">Private</option>
          </select>
        </div>
      </form>
    </div>
  );
};

// !! Main Component !!
const PostBox = ({ id }) => {
  const [body, setBody] = useState("");
  const [data, setData] = useState(null);
  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.substring(1);
      setBody(hash);
    };

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
        status: "semi-private",
        groupId: 0,
        comments: [],
        date: "",
        followers: [2, 3, 19],
      }
      fetchData(method, type, payload).then((data) => {
        /*  herf = "#postlist"; */
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
          <CreatePost onSubmit={handleSubmitPost} />
        </section>
      )}

      {body === "postlist" && (
        <section id="postlist">
          <PostList profileId={id} />
        </section>
      )}
    </>
  );
};

export default PostBox;
