import Header from "./User/Header";
import Body from "./User/Body";
import "./MainBox.css";
const MainBox = ({ user }) => {
    return (
        <div className="main-box">
            <Header user={user} />
            <Body user={user} />
        </div>
    );

}

export default MainBox;