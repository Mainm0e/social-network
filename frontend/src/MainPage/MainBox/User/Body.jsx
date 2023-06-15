import "./user.css";
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const Body = (id) => {
  const postid = id.user;
  return (
    <div className="main_body">
      {/* Conditional rendering based on the body state */}
      <PostBox id={id.user}/>
    </div>
  );
};

export default Body;

