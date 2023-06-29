import React, { useState, useEffect } from "react";
import { getCookie } from "../../tools/cookie";
import "./UserList.css";
import { fetchData } from "../../tools/fetchData";


const UserList = ({title,id,clearBox}) => {
   
    const [data, setData] = useState(null);
    useEffect(() => {
    const method = "POST"
    const type = "profileList"
    const payload = {sessionId:getCookie("sessionId"), userId: id, request:title}
    console.log(payload)
    fetchData(method,type,payload).then((data) => {
        setData(data)
    })
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