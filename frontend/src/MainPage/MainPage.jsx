import React, { useEffect, useState } from "react";
import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from "../Common/ChatBox/ChatBox";
import { navGroupLinkData } from "./dummyData";
import { getCookie, getUserId} from "../tools/cookie";
import { fetchData } from "../tools/fetchData";

function MainPage() {

  // !!TODO!! how to get profile id that can send to mainbox that show user that we want ???

  const [data, setData] = useState(null);
  useEffect(() => {
    const method = "POST"
    const type = "profile"
    const payload = { sessionId: getCookie("sessionId"), userId: getUserId("userId"), profileId: getUserId("userId") }
    fetchData(method,type, payload).then((data) => setData(data));
  }, []);

 

  if (!data) {
    return <div>Loading...</div>;
  } else {
  return (
    <div className="main-page">
      <div className="main-page-container">
        <BoxState userData={data} navLink={navGroupLinkData}/>
      </div>
    </div>
  );
}
}

export default MainPage;

// BoxState is component that read url and send user to correct place
// first check url pathname 
// then if is looking for state value that is null or not if state have some value 
// it will sent state value to profile componen
const BoxState = ({userData, navLink}) => {
  /*  !!todo!!
      - error handler that if mainbox cant find user profile with that id what need todo */
  const url = new URL(window.location.href);
  const searchParams = new URLSearchParams(url.search);
  const state = searchParams.get("id");
  if (url.pathname === "/user") {
    if (state !== null){
      return <Profile userData={userData} profileId={state} navLink={navLink} type={"user"}/>
    } else if (state === null){
     return  <Explore userData={userData} navLink={navLink} type={"user"}/>
    }
  } else if (url.pathname === "/group"){
    if (state !== null){
    } else if (state === null){
      return <Explore userData={userData} navLink={navLink} type={"group"}/>
    }
  } else if (url.pathname === "/"){
    const id = getUserId("userId")
    return <Profile userData={userData} profileId={parseInt(id)} navLink={navLink} type={"user"}/>
  }
}



const Profile = ({userData, profileId, navLink ,type}) => {
  return (
  <>
  <LeftBox user={userData} link={navLink} />
  <MainBox profileId={parseInt(profileId)} type={type} state={"profile"}/>
  <RightBox/>
  <ChatBox/>
  </> 
  )
}

const Explore = ({userData, navLink, type}) => {
  console.log("userDat",userData,"navLink", navLink,"type",type)
  return (
    <>
      <LeftBox user={userData} link={navLink}/>
      <MainBox userId={null} type={type} state={"explore"}/>
      <RightBox/>
      <ChatBox/>
    </>
  )
}