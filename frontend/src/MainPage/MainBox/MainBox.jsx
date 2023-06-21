import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import Explore from "./explore/explore";
import "./MainBox.css";
import { getCookie } from "../../tools/cookie";
const MainBox = ({ user }) => {
  const [data, setData] = useState(null);
  const sessionId = getCookie("sessionId");
  const [boxState, setBoxState] = useState(null);


  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
          type: "profile",
          payload: { sessionId: sessionId, userId: user, profileId: user },
        }),
      });
      const responseData = await response.json();
      setData(responseData.event.payload);
    };
    fetchData();
  }, []);

  const changeBoxState = (state) => {
    if (state === "profile") {
      console.log("in changeBoxState", data)
      setBoxState(profile(data,user));
    } else if (state === "explore") {
      setBoxState(explore());
    }
  };
  
  // click button to change changeBoxState input
  const [clickState, setClickState] = useState("");
  const click = () => {
    console.log("clickState", clickState)

    if (clickState === "explore") {
      setClickState("profile");
      changeBoxState(clickState);
    } else if (clickState === "profile") {
      setClickState("explore");
      changeBoxState(clickState);
    } else {
      setClickState("explore");
      changeBoxState(clickState);
    }
  };

  //return boxState === "profile" ? <Header profile={data} /> : <Body user={user} />;
  if (data === null) {
    return <div>Loading...</div>;
  } else {
    return (
      <div className="main-box">
       <button onClick={click}>click</button>
        {boxState}
      </div>
    );
  }
};

export default MainBox;

const explore = () =>{
  return (
    <>
     <Explore /> 
     </>
  )
}

const profile = (data,id) =>{
  return (
    <>
     <Header profile={data} /> 
     <Body user={id} />
     </>
  )
}