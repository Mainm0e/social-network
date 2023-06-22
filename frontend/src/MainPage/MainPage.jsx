import React, { useEffect, useState } from "react";
import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from "../Common/ChatBox/ChatBox";
import { navGroupLinkData } from "./dummyData";
import { getCookie, getUserId} from "../tools/cookie";

// dummy data
function MainPage() {
  // get userId from cookie
  const id = getUserId("userId")
  const userId =  parseInt(id);

  // !!TODO!! how to get profile id that can send to mainbox that show user that we want ???
  const profileId = 4
  // make url = localhost:3000/
  const url = window.location.href;
  const urlSplit = url.split("/");
  const urlJoin = urlSplit.slice(0, 3).join("/");
  window.history.pushState({}, null, urlJoin);

  const [data, setData] = useState(null);
  const sessionId = getCookie("sessionId");
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ type: "profile", payload: {sessionId:sessionId, userId: userId, profileId: userId} }),
      });
      const responseData = await response.json();
      setData(responseData);
    };

    fetchData();

  }, []);

  

  if (!data) {
    return <div>Loading...</div>;
  } else{
  return (
    <div className="main-page">
      <div className="main-page-container">
        <LeftBox user={data.event.payload} link={navGroupLinkData}/>
        <MainBox  profileId={profileId} state={"explore"}/>
        <RightBox/>
        <ChatBox />
      </div>
    </div>
  );
}
}

export default MainPage;


const Profile = ({userData, profileId, navLink}) => {
  return (
  <>
  <LeftBox user={userData} link={navLink} />
  <MainBox userId={profileId} state={"profile"}/>
  <RightBox/>
  <ChatBox/>
  </> 
  )
}

const Explore = ({userData, navLink}) => {
  return (
    <>
      <LeftBox user={userData} link={navLink}/>
      <MainBox userId={null} state={"explore"}/>
      <RightBox/>
      <ChatBox/>
    </>
  )
}