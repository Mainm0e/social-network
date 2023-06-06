import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from '../Common/ChatBox/ChatBox';

// dummy data
import { user1 } from "./dummyData";


function MainPage() {
  return (
    <div className="main-page">
        <div className="main-page-container">
        <LeftBox user={user1}/>
        <MainBox user={user1}/>
        <RightBox />
        <ChatBox />
        </div>
    </div>

  );
}
export default MainPage;