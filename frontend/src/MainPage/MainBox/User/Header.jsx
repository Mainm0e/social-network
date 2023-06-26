import React, { useState } from "react";
import "../MainBox.css";
import { getCookie, getUserId } from "../../../tools/cookie";
const Header = ({profile,handleRefresh}) => {
  const [requestSent, setRequestSent] = useState(false);
  const refreshHeader = () => {
    setRequestSent(false); // Reset the request status
  };
  const user = profile;
  const checkPrivacy = () => {
    if (true) {
      return (
        <>
          <div className="birthdate">
            <label>Birthdate: </label>
            <span>{user.privateProfile.birthdate}</span>
          </div>
          <div className="email">
            <label>Email: </label>
            <span>{user.privateProfile.email}</span>
          </div>
        </>
      );
    } else {
      return <></>;
    }
  };
  const sessionId = getCookie("sessionId");
  const userId = getUserId("userId");
    const sentRequest = async () => {
       const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
          type: "followRequest",
          payload: { sessionId: sessionId, userId: parseInt(userId), followId:user.userId,response:""},
        }),
      });
      const responseData = await response.json();
      handleRefresh();
    };
  
  return (
    <div className="main_header" key={user.userId}  >
      <div className="info_img">
        <img src={user.avatar} alt="user-img" />
      </div>
      <div className="info_box">
        <div className="user_info">
          <div className="fullName">
            <label> Name: </label>
            <span>{user.firstName}</span>
            <span> </span>
            <span>{user.lastName}</span>
          </div>
          {checkPrivacy()}
          <div className="followers">
            <label>Followers: </label>
            <span>{user.followerNum}</span>
          </div>
          <div className="following">
            <label>Following: </label>
            <span>{user.followingNum}</span>
          </div>
          {/* follow button */}
          <Followbtn relation={user.relation} sentRequest={sentRequest} refreshHeader={refreshHeader} />
        </div>
      </div>
    </div>
  );
};

export default Header;

const Followbtn = ({  relation, sentRequest, refreshHeader }) => {

  const handleSentRequest = async () => {
    await sentRequest();
    // Trigger the refresh of the Header component
    refreshHeader();
  };
  
  if (relation === "you"){
    return <></>;
  } else if (relation === "following") {
    return (
      <div className="follow_btn">
        <button className="follow_btn"  onClick={handleSentRequest}>
          Unfollow
        </button>
      </div>
    );
  }
  else if (relation === "follow") {
    return (
      <div className="follow_btn">
        <button className="follow_btn" onClick={handleSentRequest}>
          follow
        </button>
      </div>
    );
  } else if (relation === "pending"){
    return (
      <div className="follow_btn">
        <button className="follow_btn" onClick={handleSentRequest}>
          pending
        </button>
      </div>
    )
  } else {
    return (
      <></>
    );
  }
}
