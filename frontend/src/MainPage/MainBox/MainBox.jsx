import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import "./MainBox.css";
import { getCookie } from "../../tools/cookie";
const MainBox = ({ user }) => {
  const [data, setData] = useState(null);
  const sessionId = getCookie("sessionId");

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

  useEffect(() => {
    console.log("in header", data);
  }, [data]);

  if (data === null) {
    return <div>Loading...</div>;
  } else {
    console.log("in mainbox", data)
    return (
      <div className="main-box">
        <Header profile={data} />
        <Body user={user} />
      </div>
    );
  }
};

export default MainBox;