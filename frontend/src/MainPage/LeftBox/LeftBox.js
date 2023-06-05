
import { navGroupLinkData } from "../dummyData";
import NavList from "./NavList";
import "./LeftBox.css";
const LeftBox = ({user}) => {
    return (
        <div className="left-box">
            <div className="user_box">
                <div className="img_box">
                    <img src={user.avatar} alt="user-img" />
                </div>
                <div className="user_info">
                    <div className="username">
                        <span>{user.username}</span>
                    </div>
                    <div className="logout_btn">
                        <button>logout</button>
                    </div>
                </div>
            </div>
            <NavList type={"Nav"} links={navGroupLinkData[0]} />
            <NavList type={"Nav"} links={navGroupLinkData[0]} />
            <NavList type={"Nav"} links={navGroupLinkData[0]} />
            <NavList type={"Nav"} links={navGroupLinkData[0]} />
            <NavList type={"Nav"} links={navGroupLinkData[0]} />

        </div>
    );
}

export default LeftBox;