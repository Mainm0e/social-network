import "../MainBox.css";
import { getCookie, getUserId } from "../../../tools/cookie";
const Header = ({profile}) => {
  console.log("in header p", profile)
  const user = profile;
  console.log("in header u", user)
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
      console.log(userId)
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
      console.log("expor in sentRequest", responseData)
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
          <div className="follow_button">
            <button onClick={sentRequest}>Follow</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Header;
