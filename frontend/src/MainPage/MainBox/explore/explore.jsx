import React, { useState, useEffect } from "react";
import { getCookie,getUserId } from "../../../tools/cookie";
import "./explore.css";

const Explore = ({type}) => {
    console.log("type", type)
    const sessionId = getCookie("sessionId");
    const id = getUserId("userId")
    const uId =  parseInt(id);

    const [data, setData] = useState(null);
    useEffect(() => {
    const getExplore = async () => {
            const response = await fetch("http://localhost:8080/api", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ 
                type: "exploreUsers",
                payload: { sessionId: sessionId, userId: uId},
              }),
            });
            const responseData = await response.json();
            console.log("expor in sentRequest", responseData)
            setData(responseData.event.payload);
            };
            getExplore();
    }, []);
    
    const generateExploreList = () => {
        return data.map((user) => {
            return (
                <div className="explore_list_item" key={user.userId} userid={user.userId}>
                    <div className="explore_list_item_left">
                        <img src={user.avatar} alt="profile" />
                    </div>
                    <div className="explore_list_item_right">
                        <div className="explore_list_item_right_top">
                            <h3>{user.firstName} {user.lastName}</h3>
                        </div>
                        <div className="explore_list_item_right_bottom">
                            <button>Follow</button>
                        </div>
                    </div>
                </div>
            )
        })
    }
    if (!data) {
        return <div>Loading...</div>;
    } else {
    return (
    <div className="explore">
        <div className="explore_top">
            <h1>Explore</h1>
        </div>
        <div className="explore_body">
            <div className="explore_list">
                {generateExploreList()}
            </div>
        </div>
    </div>
    )
    }
}
export default Explore;
