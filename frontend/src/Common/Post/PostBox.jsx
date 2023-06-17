import React, { useState, useEffect } from 'react';
import { PostData } from './dummyData';
import Comment from './CommentBox';
import './Post.css';

const PostList = (cookie) => {
  // request to get post data from backend
    const getPost = () => {
/* 
      fetch('http://localhost:8000/api/post/', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Token ${cookie.token}`
        },
      })
        .then((response) => response.json())
        .then((data) => {
          console.log(data);
          PostData = data;
        }); */

        return PostData;
        };
  return (
    <div className="post_list">
      {getPost().map((post) => (
        <Post
          key={post.id}
            id={post.id}
            title={post.title}
            content={post.content}
            image={post.image}
            time={post.time}
            user={post.user}
            comments={post.comments}
        />
        ))}
    </div>
  );
};

const Post = ({ id, title, content, image, time, user, comments}) => {
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


const CreatePost = ({ cookie ,onSubmit }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [image, setImage] = useState(null);

  const handleTitleChange = (e) => {
    setTitle(e.target.value);
  };

  const handleContentChange = (e) => {
    setContent(e.target.value);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImage(file);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const postData = {
      title: title,
      content: content,
      image: image,
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
      </form>
    </div>
  );
};


// !! Main Component !!
const PostBox = (user) => {
  console.log(user)
    const [body, setBody] = useState('');
    const [data, setData] = useState(null);
    useEffect(() => {
      const fetchData = async () => {
        const response = await fetch("http://localhost:8080/api", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ event_type: "post", payload: { user_id: user.id } }),
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
  
    const handleSubmitPost = (postData) => {
      // Logic to handle the submission of the post data
      console.log('Submitted post:', postData);
      // Make API requests or perform other operations here
      // request to create post
      /* fetch('http://localhost:8000/api/post/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Token ${dummyCookie.token}`
        },
        body: JSON.stringify(postData),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log(data);
        }); */
        setBody('postlist');
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
          <PostList />
        </section>
      )}
      </>
    )

};

export default PostBox;