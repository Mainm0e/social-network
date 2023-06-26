import React from "react";
import "../MainBox.css";
import { getCookie, getUserId } from "../../../tools/cookie";
import { fetchData } from "../../../tools/fetchData";
const Header = ({profile,handleRefresh}) => {
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
      const method = "POST"
      const type = "followRequest"
      const payload ={ sessionId: sessionId, userId: parseInt(userId), followId:user.userId,response:""}
      fetchData(method,type,payload).then((data) => {console.log(data)})
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
          <Followbtn relation={user.relation} sentRequest={sentRequest}/>
        </div>
      </div>
    </div>
  );
};

export default Header;

const Followbtn = ({  relation, sentRequest }) => {

  const handleSentRequest = async () => {
    await sentRequest();
    // Trigger the refresh of the Header component
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
