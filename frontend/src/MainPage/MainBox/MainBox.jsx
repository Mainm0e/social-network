import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import Body from "./User/Body";
import "./MainBox.css";
const MainBox = ({ id }) => {
    const [data, setData] = useState(null);
    useEffect(() => {
        const fetchData = async () => {
          const response = await fetch("http://localhost:8080/api", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ event_type: "profile", payload: {sessionId:id.sessionId, userId: 1 } }),
          });
          const responseData = await response.json();
          setData(responseData);
          console.log(responseData)
        };
        fetchData();
      }, []);
    if (!data) {
        return <div>Loading...</div>;
    }else{
        return (
            <div className="main-box">
            <Header profile={data.event.payload} />
            <Body user={data.userId} />
        </div>
    );
}
        

}

export default MainBox;