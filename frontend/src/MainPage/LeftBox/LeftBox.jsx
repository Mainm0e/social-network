

import {logout} from "../../tools/logout";
import NavList from "./NavList";
import "./LeftBox.css";

const LeftBox = ({user,link}) => {
    const handleLogout = () => {
        logout();
    }
    const checkImage = () => {
        if (user.avatar === ''|| user.avatar === null|| user.avatar === undefined) {
            return null;
        } else {
            return  <div className="post_image"> <img src={user.avatar} alt="content" /> </div>;
        }
    };
    return (
        <div className="left-box">
            <div className="user_box">
                <div className="img_box">
                    {checkImage()}
                </div>
                <div className="user_info">
                    <div className="username">
                        <span>{user.firstName}</span>
                    </div>
                    <div className="logout_btn">
                        <button onClick={handleLogout}>logout</button>
                    </div>
                </div>
            </div>
            <NavList type={"Main"} links={link[3]} />
            <NavList type={"Post"} links={link[1]} />
            <NavList type={"Nav"} links={link[0]} />
            <NavList type={"Connection"} links={link[2]} />

        </div>
    );
}

export default LeftBox;