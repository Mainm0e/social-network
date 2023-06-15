import React, { useState } from "react";
import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from '../Common/ChatBox/ChatBox';

// dummy data
import { user1 } from "./dummyData";





function MainPage(name) {
  const [user, setUser] = useState(null);
  const currentUrl = window.location.href;
  // how to set current url to /user/fname-lname
 window.history.pushState({}, null, '/user/' + user1.fName + '-' + user1.lName);



/*   fetch('http://localhost:3001/api', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': cookie
    },
    body: JSON.stringify({event_type:'profile', payload: {user_id: 1}})
  }).then((response) => {
    return response.json();
  }).then((data) => {
    setUser(data);
  }) */

  const getUser = () => {
    setUser(user1);
  }
  return (
    <div className="main-page">
      {user === null ? getUser() : null}
        <div className="main-page-container">
        <LeftBox user={user}/>
        <MainBox  user={user}/>
        <RightBox />
        <ChatBox />
        </div>
    </div>

  );
}
export default MainPage;