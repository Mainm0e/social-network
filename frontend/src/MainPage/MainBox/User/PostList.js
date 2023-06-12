import React from 'react';
import { PostData } from './dummyData';
import './Post.css';

const PostList = () => {
    const getPost = () => {
        // Logic to get the post data
        console.log('Getting post data');
        console.log(PostData)
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
        />
        ))}
    </div>
  );
};

const Post = ({ id, title, content, image, time, user}) => {
  const checkImage = () => {
    console.log(image)
    if (image === ''|| image === null|| image === undefined) {
      return null;
    } else {
      return  <div className="post_image"> <img src={image} alt="content" /> </div>;
    }
  };
    console.log(id)
  return (
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
  );
};


export default PostList;
