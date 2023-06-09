import "./user.css"
import CreatePost from "./CreatePost"
import PostList from "./PostList"
import React, { useState } from 'react';
const Body = ({id}) => {
    const [body, setBody] = useState('createpost')
    const changeBody = () => {
        if (body === 'createpost') {
            setBody('postlist')
        } else {
            setBody('createpost')
        }
    }
    const handleSubmitPost = (postData) => {
        // Logic to handle the submission of the post data
        setBody('postlist')
        console.log('Submitted post:', postData);
        // Make API requests or perform other operations here
      };
    /* 
      <CreatePost id={id} onSubmit={handleSubmitPost}/>
            <PostList/>
             */
    return (
        <div className="main_body">
          {/* button for change createpost adn poslist */}
            <button onClick={changeBody}>change</button>
            {body === 'createpost' ? (
                <CreatePost id={id} onSubmit={handleSubmitPost}/>
            ) : (
                <PostList/>
            )}
        </div>
    )
}

export default Body