import "../MainBox.css"
const Header = (profile) => {
    const user = profile.user
    return (
        <div className="main_header">
            <div className="info_img">
                <img src={user.avatar} alt="user-img" />
            </div>
            <div className="info_box">
                <div className="user_info">
                    <div className="fullName">
                      <label> Name: </label><span>{user.fName}</span><span> </span><span>{user.lName}</span>
                    </div>
                    <div className="birthdate">
                        <label>Birthdate: </label> 
                        <span>{user.birthdate}</span>
                    </div>
                    <div className="email">
                        <label>Email: </label> 
                        <span>{user.email}</span>
                    </div>
                    <div className="followers">
                        <label>Followers: </label> 
                        <span>{user.followers}</span>
                    </div>
                    <div className="following">
                        <label>Following: </label> 
                        <span>{user.following}</span>
                    </div>
                </div>
                <div className="info_bio">
                    <span>{user.bio}</span>
                </div>
            </div>
        </div>
    )
}

export default Header