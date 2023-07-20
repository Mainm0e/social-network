import "../MainBox.css"
import PostBox from "../../../Common/Post/PostBox";

import React from 'react';

const GroupBody = ({id}) => {
  return (
    <div className="main_body">
      <PostBox id={id} from={"group"}/>
    </div>
  );
};

export default GroupBody;
