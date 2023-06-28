import React, { useState } from "react";
import {logout} from "../../tools/logout";
import NavList from "./NavList";
import Notification from "../../Common/Notification/Notification";
import "./LeftBox.css";
import { useEffect } from "react"; 
import { fetchData } from "../../tools/fetchData";
import { getUserId, getCookie } from "../../tools/cookie";

const LeftBox = ({user,link}) => {
    const handleLogout = () => {
        logout();
    }
    const [notificationData, setNotificationData] = useState([]);
    useEffect(() => {
        const userId = getUserId("userId");
        const sessionId = getCookie('sessionId');
        const method = 'POST';
        const type = "requestNotif";

        const payload = {
            userId: parseInt(userId),
            sessionId:sessionId,
        };
        fetchData(method, type, payload).then((data) => {
            console.log(data);
            setNotificationData(data);
        }   
        );
    }, [user]);

    const checkImage = () => {
        if (user.avatar === ''|| user.avatar === null|| user.avatar === undefined) {
            return null;
        } else {
            return  <div className="post_image"> <img src={user.avatar} alt="content" /> </div>;
        }
    };
    return (
        <div className="left-box">
            <div className="user_box">
                <div className="img_box">
                    {checkImage()}
                </div>
                <div className="user_info">
                    <div className="username">
                        <span>{user.firstName}</span>
                    </div>
                    <div className="logout_btn">
                        <button onClick={handleLogout}>logout</button>
                    </div>
                </div>
            </div>
            <NavList type={"Main"} links={link[3]} />
            <NavList type={"Post"} links={link[1]} />
            <NavList type={"Nav"} links={link[0]} />
            <NavList type={"Connection"} links={link[2]} />
            <Notification data={notificationData}/>
        </div>
    );
}

export default LeftBox;