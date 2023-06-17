import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import "./MainBox.css";
import { getCookie } from "../../tools/cookie";
const MainBox = (user) => {
    const [data, setData] = useState(null);
    const sessionId = getCookie("sessionId");
    useEffect(() => {
        const fetchData = async () => {
          const response = await fetch("http://localhost:8080/api", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ event_type: "profile", payload: {sessionId:sessionId, userId: user.id } }),
          });
          const responseData = await response.json();
          setData(responseData);

        };
        fetchData();
      }, []);
    if (!data) {
        return <div>Loading...</div>;
    }else{
        return (
            <div className="main-box">
            <Header profile={data.event.payload} />
            <Body id={user.id} />
        </div>
    );
}
        

}

export default MainBox;