import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import Explore from "./explore/explore";
import "./MainBox.css";
import { getCookie, getUserId } from "../../tools/cookie";



const MainBox = ({ profileId, type ,state}) => {
  console.log("profileId", profileId, "state", state, "type", type)
  const sessionId = getCookie("sessionId");
  const userId = getUserId("userId")

  if (state === "explore"){
    if (type ==="user"){
      return (
        <Explore type={"exploreUsers"}/> 
        )
      } else if (type === "group"){
        return (
          <Explore type={"exploreGroups"}/>
        )
      } else {
        return (
          <Explore type={"exploreUsers"}/> 
          )
      }
  } else if (state === "profile"){
    if (type === "user"){
      return(
        <Profile sessionId={sessionId} userId={userId} profileId={profileId}/>
        )
      } else if (type === "group"){
        console.log( "im group explore")
        return(
          <>
          <p>hello</p>
          </>
        )
      }
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