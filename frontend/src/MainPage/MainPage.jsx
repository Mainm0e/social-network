import React, { useEffect, useState } from "react";
import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from "../Common/ChatBox/ChatBox";

// dummy data
function MainPage() {
  const [data, setData] = useState(null);


  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("http://localhost:8080/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ event_type: "profile", payload: { user_id: 1 } }),
      });
      const responseData = await response.json();
      setData(responseData);
    };

    fetchData();
  }, []);
  if (!data) {
    return <div>Loading...</div>;
  } else if (data.success){
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
