import React from "react";
import "../MainBox.css";
import { getCookie, getUserId } from "../../../tools/cookie";
import { fetchData } from "../../../tools/fetchData";
import { link_followers, link_following } from "../../../tools/link";
const Header = ({profile,handleRefresh}) => {
  const user = profile;
  const checkPrivacy = () => {
    if (true) {
      return (
        <>
          <div className="birthdate info">
            <label>Birthdate: </label>
            <span>{user.privateProfile.birthdate}</span>
          </div>
          <div className="email info">
            <label>Email: </label>
            <span>{user.privateProfile.email}</span>
          </div>
        </>
      );
    } else {
      return <></>;
    }
  };
    const followRequest = async () => {
      const method = "POST"
      const type = "followRequest"
      const payload ={ sessionId: getCookie("sessionId"), senderId: getUserId("userId"), receiverId:user.userId}
      fetchData(method,type,payload).then((data) => {
        console.log(data)})

      /* add delay */
      setTimeout(() => {
        handleRefresh();
      }, 100);
    };
  
    const changePrivacy = async () => {
      const method = "POST"
      const type = "updatePrivacy"
      const payload ={ sessionId: getCookie("sessionId"), userId:getUserId("userId"),privacy:user.privacy}
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
          <div className="fullName info">
            <label> Name: </label>
            <span id="profile-firstName">{user.firstName}</span>
            <span> </span>
            <span>{user.lastName}</span>
          </div>
          {checkPrivacy()}
          <div className="followers info">
            {user.followerNum > 0 && <label onClick={link_followers} >Followers: </label>||<label>Followers: </label>}
            <span>{user.followerNum}</span>
          </div>
          <div className="following info">
            {user.followingNum > 0 && <label onClick={link_following} >Following: </label>||<label>Following: </label>}
            <span>{user.followingNum}</span>
          </div>
          {/* follow button */}
          <Followbtn relation={user.relation} privacy={user.privacy} followRequest={followRequest} changePrivacy={changePrivacy} />
        </div>
      </div>
    </div>
  );
};

export default Header;

const Followbtn = ({  relation, privacy ,followRequest , changePrivacy}) => {

  const handleSentRequest = async () => {
    if (relation === "you"){
      await changePrivacy();
    } else {
      await followRequest();
    }
    // Trigger the refresh of the Header component
  };
  
  if (relation === "you"){
    if (privacy === "private"){
      return (
        <div className="follow_btn">
          <label>Status: </label>
          <button className="follow_btn hover">
            private
          </button>
          <button className="follow_btn" onClick={handleSentRequest} style={{ cursor: 'pointer' }}>
            public
          </button>
        </div>
      )
    } else if (privacy === "public"){
      return (
        <div className="follow_btn">
          <label>Status: </label>
          <button className="follow_btn hover" >
            public
          </button>
          <button className="follow_btn" onClick={handleSentRequest} style={{ cursor: 'pointer' }}>
            private
          </button>
        </div>
      )
    } else {
      return (
        <></>
      )
    }
   
  } else if (relation === "following") {
    return (
      <div className="follow_btn">
        <button className="follow_btn"  onClick={handleSentRequest} style={{ cursor: 'pointer' }}>
          Unfollow
        </button>
      </div>
    );
  }
  else if (relation === "follow") {
    return (
      <div className="follow_btn">
        <button className="follow_btn" onClick={handleSentRequest} style={{ cursor: 'pointer' }}>
          follow
        </button>
      </div>
    );
  } else if (relation === "pending"){
    return (
      <div className="follow_btn">
        <button className="follow_btn" onClick={handleSentRequest} style={{ cursor: 'pointer' }}>
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
