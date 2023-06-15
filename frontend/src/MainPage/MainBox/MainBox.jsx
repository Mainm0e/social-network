import Header from "./User/Header";
import Body from "./User/Body";
import "./MainBox.css";
const MainBox = ({ user }) => {
    console.log(user.userId)
    return (
        <div className="main-box">
            <Header user={user} />
            <Body user={user.userId} />
        </div>
    );

}

export default MainBox;