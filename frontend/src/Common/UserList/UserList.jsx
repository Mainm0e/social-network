import React, { useState, useEffect } from "react";
import { getCookie } from "../../tools/cookie";
import "./UserList.css";


const UserList = ({title,id,clearBox}) => {
   
    const [data, setData] = useState(null);
    useEffect(() => {
    const getUserList = async () => {
        const response = await fetch("http://localhost:8080/api", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ type: "profileList", payload: {sessionId:getCookie("sessionId"), userId: id, request:title}}),
        });
        const responseData = await response.json();
        setData(responseData.event.payload);
    }
    getUserList();
    }, []);
    if (data===null) {
        return <div>Loading...</div>;
    }else{
    return (
        <div className="user-list">
            <div className="user-list-container">
                <div className="user-list-header">
                    <h2>{title}</h2>
                </div>
                <div className="user-list-body">
                    <ul>
                        {data.map((user) => (
                            <li key={user.userId}>
                                <div className="user-list-item">
                                    <div className="user-list-item-left">
                                        <img src={user.avatar} alt="user" />
                                    </div>
                                    <div className="user-list-item-right">
                                    <span>{user.firstName}</span><span> </span><span>{user.lastName}</span>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </div>
                <div className="user-list-footer">
                    <button onClick={clearBox}>Close</button>
                </div>
            </div>
        </div>
    );
    };
}

export default UserList;