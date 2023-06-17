
import { navGroupLinkData } from "../dummyData";
import {logout} from "../../tools/logout";
import NavList from "./NavList";
import "./LeftBox.css";

const LeftBox = ({user}) => {
    const handleLogout = () => {
        logout();
    }
    return (
        <div className="left-box">
            <div className="user_box">
                <div className="img_box">
                    <img src={user.avatar} alt="user-img" />
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
            <NavList type={"Post"} links={navGroupLinkData[1]} />
            <NavList type={"Nav"} links={navGroupLinkData[0]} />

        </div>
    );
}

export default LeftBox;