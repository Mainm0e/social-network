import React, { useEffect, useState } from "react";
import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from "../Common/ChatBox/ChatBox";

// dummy data
function MainPage() {
  const [data, setData] = useState(null);
  //get cookie from browser
  const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2){
      return parts[1]
    }
  };
  const sessionId = getCookie("sessionId");
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ event_type: "profileList", payload: {sessionId:sessionId, userId: 1, listName:"followings"} }),
      });
      const responseData = await response.json();
      setData(responseData);
      console.log(responseData)
    };

    fetchData();
  }, []);
  if (!data) {
    return <div>Loading...</div>;
  } else{
  return (
    <div className="main-page">
      <div className="main-page-container">
        <LeftBox user={data.event.payload} />
        <MainBox user={data.event.payload} />
        <RightBox />
        <ChatBox />
      </div>
    </div>
  );
}
}

export default MainPage;
