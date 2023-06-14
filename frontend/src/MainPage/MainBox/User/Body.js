import "./user.css";
import CreatePost from "./CreatePost";
import PostList from "./PostList";
import React, { useState, useEffect } from 'react';

const Body = ({ id }) => {
  const [body, setBody] = useState('');

  useEffect(() => {
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
    setBody('postlist');
    console.log('Submitted post:', postData);
    // Make API requests or perform other operations here
  };

  return (
    <div className="main_body">
      {/* Conditional rendering based on the body state */}
      {body === 'createpost' && (
        <section id="createpost">
          <CreatePost id={id} onSubmit={handleSubmitPost}/>
        </section>
      )}

      {body === 'postlist' && (
        <section id="postlist">
          <PostList/>
        </section>
      )}
    </div>
  );
};

export default Body;

