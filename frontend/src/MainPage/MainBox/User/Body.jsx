import "./user.css";
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const Body = () => {

  return (
    <div className="main_body">
      {/* Conditional rendering based on the body state */}
      <PostBox/>
    </div>
  );
};

export default Body;

