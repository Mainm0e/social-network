
import "./RightBox.css";
import UserList from "../../Common/UserList/UserList";

const RightBox = ({profileId}) => {
    //find section to render
    // #follower or #following
    return (
        <div className="right-box">
            <UserList title={"followings"} id={profileId} />
        </div>
    );
}

export default RightBox;