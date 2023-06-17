import "../MainBox.css";
const Header = (profile) => {
  const user = profile.profile;
  console.log("user",user)
    /* console.log(user.privateProfile.followers, user.privateProfile.following); */
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
  return (
    <div className="main_header" key={user.id}  >
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
        </div>
      </div>
    </div>
  );
};

export default Header;
