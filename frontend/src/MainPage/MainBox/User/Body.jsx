import "../MainBox.css"
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const Body = ({id}) => {
  return (
    <div className="main_body">
      <PostBox id={id}  from={"profile"}/>
    </div>
  );
};

export default Body;

