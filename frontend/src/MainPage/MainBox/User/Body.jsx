import "./user.css";
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const Body = (user) => {
  const postid = user.id;
  return (
    <div className="main_body">
      {/* Conditional rendering based on the body state */}
      <PostBox id={postid}/>
    </div>
  );
};

export default Body;

