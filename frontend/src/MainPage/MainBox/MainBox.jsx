import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import Explore from "./explore/explore";
import "./MainBox.css";
import { getCookie, getUserId } from "../../tools/cookie";



const MainBox = ({ profileId, state}) => {
  const sessionId = getCookie("sessionId");
  const userId = getUserId("userId")

  if (state === "explore"){
    return (
      <Explore type={"user"}/> 
    )
  } else if (state === "profile"){
    return(
      <Profile sessionId={sessionId} userId={userId} profileId={profileId}/>
    )
};
}


export default MainBox;


const Profile = ({sessionId, userId,profileId}) =>{
   const [data, setData] = useState(null);
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
          type: "profile",
          payload: { sessionId: sessionId, userId: parseInt(userId), profileId: profileId },
        }),
      });
      const responseData = await response.json();
      setData(responseData.event.payload);
    };
    fetchData();
  }, []);

  
  if (data === null) {
    return <div>Loading...</div>;
  } else {
    return (
      <div className="main-box">
        <Header profile={data}/>
        <Body user={profileId}/>
      </div>
    );
  }
}