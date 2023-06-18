import React, { useState  } from 'react';
import "./Post.css"
const CreatePost = ({ id,onSubmit }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [image, setImage] = useState(null);
  const [privacy, setPrivacy] = useState('public');

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

  const handlePrivacyChange = (e) => { 
    setPrivacy(e.target.value);
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

export default CreatePost;
