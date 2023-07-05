import "../MainBox.css"
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const Body = (user) => {
  const postid = user.user;
  return (
    <div className="main_body">
      <PostBox id={postid}/>
    </div>
  );
};

export default Body;

