import React, { useState, useEffect } from 'react';
import Comment from './CommentBox';
import './Post.css';
import { getCookie } from '../../tools/cookie';
import { checkPostData } from '../../tools/checkdata';

const PostList = ({id}) => {
  const [postData, setPostData] = useState(null);
  useEffect(() => {
    const getPost = async () => {
      const sessionId = getCookie('sessionId');
      const response = await fetch('http://localhost:8080/api', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({type: 'GetPost', payload: { sessionId: sessionId, userId: parseInt(id), postId : 5 } }),
      });
      const responseData = await response.json();
      setPostData(responseData.event.payload);
      console.log("in PostList",responseData)
    };
    getPost();
  }, []);
  const createPost = () => {
    if (postData !== null) {
         return (
           <Post
             key={postData.userId}
             id={postData.postId}
             title={postData.title}
             content={postData.content}
             image={postData.image}
             time={postData.date}
             user={postData.userId}
             comments={postData.comments}
           />
         ); 
       }
     };
   
  if (!postData) {
    return <div>Loading...</div>;
  } else  {
    return (
      <div className="post_list">
        {createPost()}
      </div>
    );
  }
  
};
const Post = ({ id, title, content, image, time, user, comments}) => {
  console.log("in post", comments)
  const checkImage = () => {
    console.log(image)
    if (image === ''|| image === null|| image === undefined) {
      return null;
    } else {
      return  <div className="post_image"> <img src={image} alt="content" /> </div>;
    }
  };
  return (
    <>
    <div className="post">
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
            <p>{user.username}</p>
          </div>
        </div>
      </div>
      <div className="post_body">
      <p className="content">{content}</p>
      </div>
      {/* button for comment and create comment */}
    </div>
      <Comment id={id} comments={comments} />
      </>
  );
};




const CreatePost = ({ onSubmit }) => {
 const [title, setTitle] = useState('');
const [content, setContent] = useState('');
const [privacy, setPrivacy] = useState('public');
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
  setTitle('');
  setContent('');
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
        <input
          type="file"
          accept="image/*"
          onChange={handleImageChange}
        />
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
);};



// !! Main Component !!
const PostBox = ({id}) => {
    const [body, setBody] = useState('');
    const [data, setData] = useState(null);
    useEffect(() => {
      // fetch data from backend
      // to get the post list
      const fetchData = async () => {
        const response = await fetch("http://localhost:8080/api", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ event_type: "post", payload: { user_id: id} }),
        });
        const responseData = await response.json();
        setData(responseData);
      };
  
      fetchData();
      
      const handleHashChange = () => {
        const hash = window.location.hash.substring(1);
        setBody(hash);
      };
  
      // Listen for hash changes in the URL
      window.addEventListener('hashchange', handleHashChange);
      handleHashChange(); // Initialize the body state based on the current hash
  
      // Cleanup the event listener on unmount
      return () => {
        window.removeEventListener('hashchange', handleHashChange);
      };

      
    }, []);


    // submit post
    const handleSubmitPost = (postData) => {
      // Logic to handle the submission of the post data
    const check = checkPostData(postData);
    if
     (check.status === true ){
      const sessionId = getCookie("sessionId");
      // Make API requests or perform other operations here
      // request to create post
      const createPost = async () => {
        console.log("image",postData.image)
        const response = await fetch("http://localhost:8080/api", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
    
          body: JSON.stringify({ type: "createPost", payload: {sessionId:sessionId,postId:0,userId: id, title: postData.title, content: postData.content, image: postData.image, status: "semi-private", groupId: 0, comments: [],date:"",followers:[2,3,19]}}),
        });
        const responseData = await response.json();
      }
      createPost();
    } else {
      alert(check.message)
    }
    };

    return (
        <>
        {body === 'createpost' && (
        <section id="createpost">
          <CreatePost onSubmit={handleSubmitPost}/>
        </section>
        )}

        {body === 'postlist' && (
        <section id="postlist">
          <PostList id={id} />
        </section>
      )}
      </>
    )

};

export default PostBox;