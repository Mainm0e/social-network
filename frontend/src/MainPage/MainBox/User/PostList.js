import React from 'react';
import { PostData } from './dummyData';

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
        />
        ))}
    </div>
  );
};

const Post = ({ id, title, content, image, time }) => {
    console.log(id)
  return (
    <div className="post">
      <h2>{title}</h2>
      <p>{content}</p>
      {image && <img src={image} alt="Post" />}
      <p>{time}</p>
    </div>
  );
};

export default PostList;
